package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Name            string
	Conn            *websocket.Conn
	Out             chan *Message
	LeaveChatroomCh chan *User
}

func NewUser(name string, conn *websocket.Conn, leave chan *User) *User {
	user := &User{
		Name:            name,
		Conn:            conn,
		Out:             make(chan *Message),
		LeaveChatroomCh: leave,
	}

	go user.Send()
	return user
}

func (u *User) Send() {
	for {
		select {
		case msg, ok := <-u.Out:
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Out channel closed.
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := u.Conn.WriteJSON(msg)
			if err != nil {
				log.WithError(err).Error("Cannot output websocket message")
				u.LeaveChatroomCh <- u
				return
			}
		}
	}
}

type UsersDao interface {
	CreateUser(name string) (affected int, err error)
	ExistsUser(name string) (exists bool, err error)
	DeleteAllUsers() (err error)
}

type UsersSQLDao struct {
	db *sql.DB
}

func NewUsersSQLDao(db *sql.DB) *UsersSQLDao {
	return &UsersSQLDao{
		db: db,
	}
}

func (c *UsersSQLDao) CreateUser(name string) (affected int, err error) {
	query := fmt.Sprintf(`
		INSERT IGNORE INTO %s (name)
		VALUES (?)`, usersTable)
	log.Debug(query)

	result, err := c.db.Exec(query, name)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (c *UsersSQLDao) DeleteAllUsers() (err error) {
	query := fmt.Sprintf(`
		DELETE
		FROM %s`, usersTable)
	log.Debug(query)

	_, err = c.db.Exec(query)
	return err
}

func (c *UsersSQLDao) ExistsUser(name string) (exists bool, err error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE name = ?
		)`, usersTable)
	log.Debug(query)
	err = c.db.QueryRow(query, name).Scan(&exists)
	return
}
