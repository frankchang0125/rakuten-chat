package services

import (
	"database/sql"

	"rakuten.co.jp/chatroom/models"
)

var usersDao models.UsersDao

func InitUsersService(db *sql.DB) {
	usersDao = models.NewUsersSQLDao(db)
}

func CreateUser(userName string) error {
	affects, err := usersDao.CreateUser(userName)
	if err != nil {
		return err
	}

	if affects == 0 {
		return ErrUserAlreadyExists
	}
	return nil
}

func DeleteAllUsers() error {
	return usersDao.DeleteAllUsers();
}
