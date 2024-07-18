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

	// Initialize Redis
	common.InitRedis()

	// Create a new context
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout", middlewares.AuthenticationMiddleware(), handlers.Logout)

	// Protected todo routes
	todoRoutes := r.Group("/todo")
	todoRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		todoRoutes.POST("/", handlers.CreateTodo)
		todoRoutes.GET("/:id", handlers.ReadTodo)
		todoRoutes.PUT("/:id", handlers.UpdateTodo)
		todoRoutes.DELETE("/:id", handlers.DeleteTodo)
	}

	// Protected user routes
	userRoutes := r.Group("/user")
	userRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		userRoutes.GET("/", handlers.GetUser)
		userRoutes.DELETE("/", handlers.DeleteUser)
		userRoutes.PUT("/", handlers.ChangePassword)
	}

	r.GET("redis-keys", func(c *gin.Context) {
		keys, err := common.RedisClient.Keys(c, "*").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"keys": keys})
	})

	zap.L().Info(
		"Server running",
		zap.String("port", env.ApplicationPort),
	)
	r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
}
