package common

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

func InitRedis(addr string) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
		panic(err)
	}

	zap.L().Info("Connected to Redis")

	RedisClient = client
}

func SetRedisClient(client *redis.Client) {
	RedisClient = client
}
