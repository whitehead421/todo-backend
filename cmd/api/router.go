package main

import (
	"github.com/gin-gonic/gin"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/pkg/middlewares"
)

func InitializeRoutes() *gin.Engine {
	var todoHandler handlers.TodoHandler = handlers.NewTodoHandler()
	var userHandler handlers.UserHandler = handlers.NewUserHandler()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

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
