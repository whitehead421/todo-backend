package common

import (
	"github.com/redis/go-redis/v9"
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

func SetRedisClient(client *redis.Client) {
	RedisClient = client
}
