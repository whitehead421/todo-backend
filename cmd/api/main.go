package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/common"
)

func main() {
	env := common.GetEnvironmentVariables()

	// Initialize logger
	logger := common.InitLogger()
	defer logger.Sync() // flushes buffer, if any

	// Connect to database
	common.ConnectDatabase(env.DatabaseDsn)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/todo", handlers.CreateTodo)
	r.GET("/todo/:id", handlers.ReadTodo)
	r.PUT("/todo/:id", handlers.UpdateTodo)
	r.DELETE("/todo/:id", handlers.DeleteTodo)

	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	zap.L().Info(
		"Server running",
		zap.String("port", env.ApplicationPort),
	)
	r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
}
