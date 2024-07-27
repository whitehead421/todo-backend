package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	Logout(context *gin.Context)
	Authorize(context *gin.Context)
	Verify(context *gin.Context)
}

type authHandler struct {
	validate *validator.Validate
}

func NewAuthHandler() AuthHandler {
	return &authHandler{
		validate: validator.New(),
	}
}

func (h *authHandler) Register(context *gin.Context) {
	var registerRequest models.RegisterRequest

	if err := context.ShouldBindJSON(&registerRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(registerRequest); err != nil {
		zap.L().Error("Validation error",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	result := common.DB.Where("email = ?", registerRequest.Email).First(&user)
	if result.Error == nil {
		zap.L().Error("This email is already registered",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
		return
	}

	user = entities.User{
		Email:       registerRequest.Email,
		Name:        registerRequest.Name,
		Password:    common.HashPassword(registerRequest.Password),
		Verified:    false,
		VerifyToken: common.GenerateUUID(),
	}

	result = common.DB.Create(&user)
	if result.Error != nil {
		zap.L().Error("Failed to create user",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	env := common.GetEnvironmentVariables()
	writer := common.NewKafkaWriter(env)
	defer writer.Close()

	err := writer.WriteMessages(context,
		kafka.Message{
			Key:   []byte(user.Email),
			Value: []byte(user.VerifyToken),
		},
	)
	if err != nil {
		zap.L().Error("Failed to write message to Kafka",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email", "message": err.Error()})
		context.Abort()
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

func (h *authHandler) Login(context *gin.Context) {
	var loginRequest models.LoginRequest

	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(loginRequest); err != nil {
		zap.L().Error("Validation error",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	result := common.DB.Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		zap.L().Error("Failed to find user",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !common.CheckPasswordHash(loginRequest.Password, user.Password) {
		zap.L().Error("Invalid password",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.Verified {
		zap.L().Error("Account is not verified",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not verified"})
		return
	}

	token, err := common.CreateToken(user.ID)
	if err != nil {
		zap.L().Error("Failed to create token",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = common.RedisClient.Set(context, token, "token", time.Hour).Err()
	if err != nil {
		zap.L().Error("Failed to set token to redis",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
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

func (h *authHandler) Logout(context *gin.Context) {
	userID, _ := context.Get("userID")
	authHeader := context.GetHeader("Authorization")

	if authHeader == "" {
		zap.L().Error("Authorization header is missing")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		context.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		zap.L().Error("Authorization header format must be Bearer {token}")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		context.Abort()
		return
	}

	err := common.RedisClient.Del(context, tokenString).Err()
	if err != nil {
		zap.L().Error("Failed to delete token from redis",
			zap.Error(err),
		)
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token from redis"})
		context.Abort()
		return
	}

	zap.L().Info("User logged out successfully",
		zap.Uint64("user ID", userID.(uint64)),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(200, gin.H{"message": "Successfully logged out"})
}

func (h *authHandler) Authorize(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		zap.L().Error("Authorization header is missing")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		context.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		zap.L().Error("Authorization header format must be Bearer {token}")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		context.Abort()
		return
	}

	id, err := common.ValidateToken(tokenString)
	if err != nil {
		zap.L().Error("Token is not valid anymore.", zap.Error(err))
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid anymore."})
		context.Abort()
		return
	}

	_, err = common.RedisClient.Get(context, tokenString).Result()
	if err != nil {
		zap.L().Error("Token not found in redis", zap.Error(err))
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid anymore"})
		context.Abort()
		return
	}

	var user entities.User
	result := common.DB.First(&user, id)
	if result.Error != nil {
		zap.L().Error("User not found for authentication middleware", zap.Error(result.Error))
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid anymore or user does not exist"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Authorized", "user_id": id})
}

func (h *authHandler) Verify(context *gin.Context) {
	verifyToken := context.Query("token")
	if verifyToken == "" {
		zap.L().Error("Verify token is missing",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Verify token is missing"})
		return
	}

	var user entities.User
	result := common.DB.Where("verify_token = ?", verifyToken).First(&user)
	if result.Error != nil {
		zap.L().Error("Failed to find user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Verified = true
	result = common.DB.Save(&user)
	if result.Error != nil {
		zap.L().Error("Failed to update user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	zap.L().Info("User verified successfully",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}
