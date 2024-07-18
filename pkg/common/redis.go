package common

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(addr string) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	RedisClient = client
}
