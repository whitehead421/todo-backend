package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/models"
	"go.uber.org/zap"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	env := common.GetEnvironmentVariables()
	path := fmt.Sprintf("http://auth:%s/authorize", env.AuthPort)

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		req, _ := http.NewRequest("POST", path, nil)
		req.Header.Set("Authorization", token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			zap.L().Error("Failed to authorize", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "status": resp.StatusCode})
			c.Abort()
			return
		}

		var authResponse models.AuthorizeResponse
		err = json.NewDecoder(resp.Body).Decode(&authResponse)
		if err != nil {
			zap.L().Error("Failed to parse response", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Set("userID", authResponse.UserId)
		c.Next()
	}
}
