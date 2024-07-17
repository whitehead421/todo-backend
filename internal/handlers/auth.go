package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func Register(context *gin.Context) {
	var registerRequest models.RegisterRequest

	if err := context.ShouldBindJSON(&registerRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(registerRequest.Password)
	if err != nil {
		zap.L().Error("Failed to hash password",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{
		Email:    registerRequest.Email,
		Name:     registerRequest.Name,
		Password: hashedPassword,
	}

	result := common.DB.Create(&user)
	if result.Error != nil {
		zap.L().Error("Failed to create todo",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	registerResponse := models.RegisterResponse{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	zap.L().Info("User created successfully",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusCreated, registerResponse)
}

func Login(context *gin.Context) {
	var loginRequest models.LoginRequest

	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	result := common.DB.Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		zap.L().Error("Failed to find user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(loginRequest.Password, user.Password) {
		zap.L().Error("Invalid password",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := createToken(user.ID)
	if err != nil {
		zap.L().Error("Failed to create token",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginResponse := models.LoginResponse{
		Token:  token,
		UserId: user.ID,
	}

	zap.L().Info("User logged in successfully",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, loginResponse)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createToken(id uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Hour).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	}

	return "", err
}
