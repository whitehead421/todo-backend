package common

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	RedisClient = client
}
