package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/whitehead421/todo-backend/pkg/common"
)

func main() {
	env := common.GetEnvironmentVariables()

	// Initialize logger
	logger := common.InitLogger()
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if err != nil {
			zap.L().Error("Failed to sync logger", zap.Error(err))
		}
	}()

	// Connect to database
	common.ConnectDatabase(env.DatabaseDsn)

	// Initialize Redis
	common.InitRedis(env.RedisAddr)

	// Initialize routes
	r := InitializeRoutes()

	// Kafka Reader
	kafkaReader := common.NewKafkaReader(env)
	defer kafkaReader.Close()

	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go common.SendActivationMail(kafkaReader, ctx)

	zap.L().Info(
		"Notification service is running",
		zap.String("port", env.NotificationPort),
	)
	err := r.Run(fmt.Sprintf(":%s", env.NotificationPort))
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
