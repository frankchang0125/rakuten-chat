package services

import (
	"database/sql"

	"github.com/go-redis/redis"
	"gopkg.in/olivere/elastic.v6"
)

var db *sql.DB
var esClient *elastic.Client
var redisClient *redis.Client

func Init(_db *sql.DB, _esClient *elastic.Client, _redisClient *redis.Client) {
	// SQL database
	db = _db
	esClient = _esClient
	redisClient = _redisClient
	InitChatroomService(db)
    InitUsersService(db)
}
