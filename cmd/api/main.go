package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/middlewares"
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

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected todo routes
	todoRoutes := r.Group("/todo")
	todoRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		todoRoutes.POST("/", handlers.CreateTodo)
		todoRoutes.GET("/:id", handlers.ReadTodo)
		todoRoutes.PUT("/:id", handlers.UpdateTodo)
		todoRoutes.DELETE("/:id", handlers.DeleteTodo)
	}

	zap.L().Info(
		"Server running",
		zap.String("port", env.ApplicationPort),
	)
	r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
}
