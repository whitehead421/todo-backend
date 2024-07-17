package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"go.uber.org/zap"
)

var secretKey = []byte(common.GetEnvironmentVariables().JwtSecret)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			zap.L().Error("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			zap.L().Error("Authorization header format must be Bearer {token}")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		id, err := ValidateToken(tokenString)
		if err != nil {
			zap.L().Error("Invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user still exists
		var user entities.User
		result := common.DB.First(&user, id)
		if result.Error != nil {
			zap.L().Error("User not found for authentication middleware", zap.Error(result.Error))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid anymore or user does not exist"})
			c.Abort()
			return
		}

		// Set userID in context
		c.Set("userID", id)
		c.Next()
	}
}

func ValidateToken(tokenString string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat, ok := claims["id"].(float64)
		if !ok {
			return 0, err
		}
		return uint64(idFloat), nil
	}

	return 0, err
}
