package redis

import (
	"fmt"

	_redis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var redisClient *_redis.Client

// InitClient initialize Redis client.
// Note: This not thread-safe and should be called at server startup.
func InitClient() *_redis.Client {
	if redisClient == nil {
		host := viper.GetString("redis.host")
		port := viper.GetInt("redis.port")
		addr := fmt.Sprintf("%s:%d", host, port)

		log.WithFields(log.Fields{
			"addr": addr,
		}).Info("Initialize Redis client")

		redisClient = _redis.NewClient(&_redis.Options{
			Addr: addr,
		})
	}
	
	return redisClient
}
