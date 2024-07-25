package main

import (
	"github.com/gin-gonic/gin"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/middlewares"
)

func InitializeRoutes() *gin.Engine {
	var authHandler handlers.AuthHandler = handlers.NewAuthHandler()
	var todoHandler handlers.TodoHandler = handlers.NewTodoHandler()
	var userHandler handlers.UserHandler = handlers.NewUserHandler()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	router.POST("/logout", middlewares.AuthenticationMiddleware(), authHandler.Logout)

	// Protected todo routes
	todoRoutes := router.Group("/todo")
	todoRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		todoRoutes.POST("/", todoHandler.CreateTodo)
		todoRoutes.GET("/:id", todoHandler.ReadTodo)
		todoRoutes.PUT("/:id", todoHandler.UpdateTodo)
		todoRoutes.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Protected user routes
	userRoutes := router.Group("/user")
	userRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		userRoutes.GET("/", userHandler.GetUser)
		userRoutes.DELETE("/", userHandler.DeleteUser)
		userRoutes.PUT("/", userHandler.ChangePassword)
	}

	return router
}
