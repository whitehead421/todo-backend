package main

import (
	"github.com/gin-gonic/gin"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/middlewares"
)

func InitializeRoutes() *gin.Engine {
	var authHandler handlers.AuthHandler = handlers.NewAuthHandler()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/authorize", authHandler.Authorize)
	router.GET("/verify", authHandler.Verify)

	router.POST("/logout", middlewares.AuthenticationMiddleware(), authHandler.Logout)

	return router
}
