package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/whitehead421/todo-backend/docs"
)

// @title Todo API
// @version 1.0
// @description This is a simple todo API
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("", CreateTodo)
	r.GET("", ReadTodo)
	r.PUT("", UpdateTodo)
	r.DELETE("", DeleteTodo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	log.Default().Println("Server running on port 8080")
	r.Run(":8080")
}

// @Summary Create todo
// @Description Create todo
// @Produce  json
// @Param todo body entities.TodoRequest true "Todo Request"
// @Success 200 {object} entities.TodoResponse
// @Router / [post]
func CreateTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Get Todo")
}

// @Summary Get todo
// @Description Get todo
// @Produce  json
// @Router / [get]
func ReadTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Read Todo")
}

// @Summary Update todo
// @Description Update todo
// @Produce  json
// @Router / [put]
func UpdateTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Update Todo")
}

// @Summary Delete todo
// @Description Delete todo
// @Produce  json
// @Router / [delete]
func DeleteTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Delete Todo")
}
