package services

import (
	"errors"
)

// Elasticsearch
const messagesIndexName = "messages"
const messageIndexType = "_doc"

// Redis
const chatroomsSet = "chatrooms"
const chatroomBaseName = "chatroom"

var (
	ErrChatroomNotExists = errors.New("Chatroom not exists")
	ErrUserNotExists     = errors.New("User not exists")
	ErrUserAlreadyExists = errors.New("User already exists")
)
