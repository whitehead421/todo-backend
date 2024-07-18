package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	Logout(context *gin.Context)
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
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(registerRequest); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := common.HashPassword(registerRequest.Password)
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

func (h *authHandler) Login(context *gin.Context) {
	var loginRequest models.LoginRequest

	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(loginRequest); err != nil {
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

	if !common.CheckPasswordHash(loginRequest.Password, user.Password) {
		zap.L().Error("Invalid password",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := common.CreateToken(user.ID)
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

func (h *authHandler) Logout(context *gin.Context) {
	// Add the token to the blacklist
	authHeader := context.GetHeader("Authorization")

	tokenString := authHeader[len("Bearer "):]
	expiration := time.Hour // Token stays in blacklist for 1 hour

	err := common.BlacklistToken(tokenString, expiration, context)
	if err != nil {
		context.JSON(500, gin.H{"error": "Failed to blacklist token"})
		return
	}

	context.JSON(200, gin.H{"message": "Successfully logged out"})
}
