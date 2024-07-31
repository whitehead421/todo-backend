package main

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Notification service is running",
		})
	})

	return router
}
