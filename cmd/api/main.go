package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/whitehead421/todo-backend/docs"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/common"
)

// @title Todo API
// @version 1.0
// @description This is a simple todo API
func main() {
	env := common.GetEnvironmentVariables()

	// Initialize logger
	logger := common.InitLogger()
	defer logger.Sync() // flushes buffer, if any

	// Connect to database
	common.ConnectDatabase(env.DatabaseDsn)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("", handlers.CreateTodo)
	r.GET("/:id", handlers.ReadTodo)
	r.PUT("/:id", handlers.UpdateTodo)
	r.DELETE("/:id", handlers.DeleteTodo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	zap.L().Info(
		"Server running",
		zap.String("port", env.ApplicationPort),
	)
	r.Run(fmt.Sprintf(":%s", env.ApplicationPort))
}
