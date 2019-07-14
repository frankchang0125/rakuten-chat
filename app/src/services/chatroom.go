package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"rakuten.co.jp/chatroom/models"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
)

const historyMessagesSize = 100

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chatroomsDao models.ChatroomsDao
var chatrooms = map[int]*models.Chatroom{}
var chatroomLock = sync.RWMutex{}
var backlog chan *models.Message

func InitChatroomService(db *sql.DB) {
	chatroomsDao = models.NewChatroomsSQLDao(db)
	backlog = make(chan *models.Message)
	go logger()
}

func CreateChatroom() (id int, err error) {
	id, err = chatroomsDao.CreateChatroom()
	if err != nil {
		return
	}

	// Invalidate chat roms list in Redis.
	err = redisClient.Del(chatroomsSet).Err()
	if err != nil {

	}

	return id, nil
}

// GetChatroomsList returns the IDs of all chat rooms.
func GetChatroomsList() ([]int, error) {
	// Try to get chat rooms list from Redis.
	// Note: We cannot distinguish nil and empty set in Redis,
	//		 thus it would be possible to hit database with true empty set.
	//		 The solution would be to wrap EXISTS and SMEMBERS commands
	//		 in a transaction. But empty set of chat rooms be a rare case,
	//		 so we ignore it.
	{
		_ids, err := redisClient.SMembers(chatroomsSet).Result()
		if err != nil {
			log.WithError(err).Error("Fail to get chat rooms list")
		} else if len(_ids) > 0 {
			ids := make([]int, 0, len(_ids))
			for _, id := range _ids {
				id, _ := strconv.Atoi(id)
				ids = append(ids, id)
			}

			return ids, err
		}
	}

	// Chat rooms list not found in Redis, retrieve from database.
	ids, err := chatroomsDao.GetChatrooms()
	if err != nil {
		return nil, err
	}

	// Save chat rooms list to Redis.
	{
		_ids := make([]interface{}, 0, len(ids))
		for _, id := range ids {
			_ids = append(_ids, id)
		}

		err = redisClient.SAdd(chatroomsSet, _ids...).Err()
		if err != nil {
			log.WithError(err).Error("Fail to set chat rooms list")
		}
	}

	return ids, nil
}

func GetChatroomUsers(chatroomID int) ([]string, error) {
	exists, err := chatroomsDao.ExistsChatroom(chatroomID)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, ErrChatroomNotExists
	}

	// Try to get chat room users list from Redis.
	// Note: We cannot distinguish nil and empty set in Redis,
	//		 thus it would be possible to hit database with true empty set.
	//		 The solution would be to wrap EXISTS and SMEMBERS commands
	//		 in a transaction. But empty set of chat rooms be a rare case,
	//		 so we ignore it.
	chatroom := chatroomSetName(chatroomID)
	users, err := redisClient.SMembers(chatroom).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"chatroom": chatroomID,
		}).Error("Fail to get chat room users list")
	} else {
		if len(users) > 0 {
			return users, nil
		}
	}

	// Chat room users list not found in Redis, retrieve from database.
	users, err = chatroomsDao.GetChatroomUsers(chatroomID)
	if err != nil {
		return nil, err
	}

	// Save chat room users list to Redis.
	if len(users) > 0 {
		_users := make([]interface{}, 0, len(users))
		for _, user := range users {
			_users = append(_users, user)
		}

		err = redisClient.SAdd(chatroom, _users...).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"err":      err,
				"chatroom": chatroomID,
			}).Error("Fail to save chat room users list")
		}
	}

	return users, nil
}

func JoinChatroom(w http.ResponseWriter, r *http.Request,
	chatroomID int, userName string) error {
	// Check if chatroom exists.
	exists, err := chatroomsDao.ExistsChatroom(chatroomID)
	if err != nil {
		return err
	} else if !exists {
		return ErrChatroomNotExists
	}

	// Check if user exists.
	exists, err = usersDao.ExistsUser(userName)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserNotExists
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	err = chatroomsDao.JoinChatroom(chatroomID, userName)
	if err != nil {
		return err
	}

	// Add user to chat room users list in Redis cache.
	chatroom := chatroomSetName(chatroomID)
	err = redisClient.SAdd(chatroom, userName).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"chatroom": chatroomID,
			"user":     userName,
		}).Error("Fail to add user to chat room users list")
	}

	chatroomLock.Lock()
	c, ok := chatrooms[chatroomID]
	if !ok {
		// This chatroom may exists in other server,
		// create a new one and subscribe to the channel.
		pubsub := redisClient.Subscribe(chatroom)
		c = models.NewChatroom(pubsub.Channel())
		chatrooms[chatroomID] = c
	}
	chatroomLock.Unlock()

	user := models.NewUser(userName, ws, c.LeaveCh)
	c.JoinCh <- user

	// Push history messages to user.
	histories, err := getHistoryMessages(chatroomID)
	if err != nil {
		if !elastic.IsNotFound(err) {
			log.WithError(err).Error("Fail to get history messages")
			return err
		}
	} else {
		go func() {
			for i := len(histories) - 1; i >= 0; i-- {
				user.Out <- histories[i]
			}
		}()
	}

	log.WithFields(log.Fields{
		"chatroom": chatroomID,
		"user":     userName,
	}).Info("Joined chat room")

	for {
		// Read message from client.
		_, msg, err := ws.ReadMessage()
		if err != nil {
			defer func() {
				err := chatroomsDao.LeaveChatroom(chatroomID, userName)
				if err != nil {
					log.WithFields(log.Fields{
						"err":      err,
						"chatroom": chatroomID,
						"user":     userName,
					}).Error("Fail to remove user from MySQL chat room users list")
				}

				// Remove user from chat room users list in Redis cache.
				err = redisClient.SRem(chatroom, userName).Err()
				if err != nil {
					log.WithFields(log.Fields{
						"err":      err,
						"chatroom": chatroomID,
						"user":     userName,
					}).Error("Fail to remove user from Redis chat room users list")
				}

				c.LeaveCh <- user
			}()

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Info("Connection closed")
				return nil
			}

			log.WithError(err).Error("Failed to read websocket message")
			return err
		}

		msgStr := string(msg)
		log.Debug(msgStr)

		m := &models.Message{
			Chatroom:  chatroomID,
			User:      userName,
			Message:   msgStr,
			Timestamp: time.Now().UnixNano(),
		}

		// Backlog message.
		backlog <- m

		// Broadcast message to all users in the same chatroom.
		_, err = redisClient.Publish(chatroom, m).Result()
		if err != nil {
			log.WithError(err).WithField("msg", m).Error("Fail to publish message")
			return err
		}
	}
}

func getHistoryMessages(chatroomID int) ([]*models.Message, error) {
	boolQuery := elastic.NewBoolQuery()
	chatroomTerm := elastic.NewTermQuery("chatroom", chatroomID)
	boolQuery.Filter(chatroomTerm)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := esClient.Search(messagesIndexName).
		Type(messageIndexType).
		Query(boolQuery).
		Sort("timestamp", false).
		Size(historyMessagesSize).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	backlogs := make([]*models.Message, 0)

	for _, hit := range resp.Hits.Hits {
		backlog := models.Message{}

		jsonStr, err := hit.Source.MarshalJSON()
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonStr, &backlog)
		if err != nil {
			return nil, err
		}

		backlogs = append(backlogs, &backlog)
	}

	return backlogs, nil
}

func logger() {
	ctx := context.Background()
	for b := range backlog {
		b.Timestamp = time.Now().UnixNano()
		_, err := esClient.Index().
			Index(messagesIndexName).
			Type(messageIndexType).
			BodyJson(b).
			Do(ctx)
		if err != nil {
			log.WithError(err).Error("Fail to backlog message")
		}
	}
}

func chatroomSetName(chatroomID int) string {
	return fmt.Sprintf("%s:%d", chatroomBaseName, chatroomID)
}
