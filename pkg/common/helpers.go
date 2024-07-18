package common

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BlacklistToken(tokenString string, expiration time.Duration, ctx context.Context) error {
	err := RedisClient.Set(ctx, tokenString, "blacklisted", expiration).Err()
	return err
}

func IsTokenBlacklisted(tokenString string, ctx context.Context) (bool, error) {
	val, err := RedisClient.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "blacklisted", err
}
