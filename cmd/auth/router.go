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

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	router.POST("/logout", middlewares.AuthenticationMiddleware(), authHandler.Logout)

	router.POST("/authorize", authHandler.Authorize)

	return router
}
