package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Chatroom  int    `json:"chatroom"`
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// MarshalBinary implements BinaryMarshaler interface.
func (m *Message) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

type Chatroom struct {
	ID      int
	clients map[string]*User
	MsgCh   <-chan *redis.Message
	JoinCh  chan *User
	LeaveCh chan *User
}

func NewChatroom(msgCh <-chan *redis.Message) *Chatroom {
	c := &Chatroom{
		clients: map[string]*User{},
		MsgCh:   msgCh,
		JoinCh:  make(chan *User),
		LeaveCh: make(chan *User),
	}

	go c.Controller()
	return c
}

func (c *Chatroom) Controller() {
	for {
		select {
		case user := <-c.JoinCh:
			c.clients[user.Name] = user
		case user := <-c.LeaveCh:
			if user, ok := c.clients[user.Name]; ok {
				close(user.Out)
				delete(c.clients, user.Name)
			}
		case msg := <-c.MsgCh:
			log.WithFields(log.Fields{
				"msg": msg.Payload,
			}).Debug("Broadcast message")

			m := Message{}
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.WithError(err).WithField("msg", msg.Payload).
					Error("Error unmarshal message")
				continue
			}

			for _, client := range c.clients {
				client.Out <- &m
			}
		}
	}
}

type ChatroomsDao interface {
	CreateChatroom() (id int, err error)
	GetChatrooms() ([]int, error)
	GetChatroomUsers(id int) (users []string, err error)
	DeleteChatroom(id string) (err error)
	JoinChatroom(chatroomID int, userName string) (err error)
	LeaveChatroom(chatroomID int, userName string) (err error)
	ExistsChatroom(chatroomID int) (bool, error)
}

type ChatroomsSQLDao struct {
	db *sql.DB
}

func NewChatroomsSQLDao(db *sql.DB) *ChatroomsSQLDao {
	return &ChatroomsSQLDao{
		db: db,
	}
}

func (c *ChatroomsSQLDao) CreateChatroom() (id int, err error) {
	query := fmt.Sprintf(`
		INSERT INTO %s ()
		VALUES ()`, chatroomsTable)
	log.Debug(query)
	result, err := c.db.Exec(query)
	if err != nil {
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return
	}

	return int(lastInsertID), nil
}

func (c *ChatroomsSQLDao) GetChatrooms() (chatrooms []int, err error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM %s`, chatroomsTable)
	log.Debug(query)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}

	chatrooms = []int{}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		chatrooms = append(chatrooms, id)
	}

	return chatrooms, nil
}

func (c *ChatroomsSQLDao) GetChatroomUsers(id int) (users []string, err error) {
	query := fmt.Sprintf(`
		SELECT name
		FROM %s AS u
		INNER JOIN %s AS cu
		ON cu.user = u.name
		WHERE cu.chatroom = ?`, usersTable, chatroomUsersTable)
	log.Debug(query)
	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	users = []string{}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		users = append(users, name)
	}

	return users, nil
}

func (c *ChatroomsSQLDao) DeleteChatroom(id string) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`
		DELETE
		FROM %s
		WHERE %s.chatroom = ?`,
		chatroomUsersTable, chatroomUsersTable)
	log.Debug(query)
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		DELETE
		FROM %s
		WHERE %s.id = ?`,
		chatroomsTable, chatroomsTable)
	log.Debug(query)
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatroomsSQLDao) JoinChatroom(chatroomID int, userName string) (err error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (chatroom, user)
		VALUES (?, ?)`, chatroomUsersTable)
	log.Debug(query)
	_, err = c.db.Exec(query, chatroomID, userName)
	return err
}

func (c *ChatroomsSQLDao) LeaveChatroom(chatroomID int, userName string) (err error) {
	query := fmt.Sprintf(`
		DELETE
		FROM %s
		WHERE %s.chatroom = ? AND %s.user = ?`,
		chatroomUsersTable, chatroomUsersTable, chatroomUsersTable)
	log.Debug(query)
	_, err = c.db.Exec(query, chatroomID, userName)
	return err
}

func (c *ChatroomsSQLDao) ExistsChatroom(chatroomID int) (exists bool, err error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE id = ?
		)`, chatroomsTable)
	log.Debug(query)
	err = c.db.QueryRow(query, chatroomID).Scan(&exists)
	return
}
