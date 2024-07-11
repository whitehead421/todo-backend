package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("", CreateTodo)
	r.GET("", ReadTodo)
	r.PUT("", UpdateTodo)
	r.DELETE("", DeleteTodo)

	r.Run()
}

func CreateTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Get Todo")
}

func ReadTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Read Todo")
}

func UpdateTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Update Todo")
}

func DeleteTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Delete Todo")
}
