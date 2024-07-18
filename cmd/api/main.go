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
	defer logger.Sync() // flushes buffer, if any

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
	r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
}
