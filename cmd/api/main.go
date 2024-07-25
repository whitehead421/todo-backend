package main

import (
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

	zap.L().Info(
		"Server running",
		zap.String("port", env.ApplicationPort),
	)
	err := r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
