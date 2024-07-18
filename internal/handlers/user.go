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

type UserHandler interface {
	GetUser(context *gin.Context)
	DeleteUser(context *gin.Context)
	ChangePassword(context *gin.Context)
}

type userHandler struct {
	validate *validator.Validate
}

func NewUserHandler() UserHandler {
	return &userHandler{
		validate: validator.New(),
	}
}

func (h *userHandler) GetUser(context *gin.Context) {
	userID, _ := context.Get("userID")

	var user entities.User

	result := common.DB.First(&user, userID)
	if result.Error != nil {
		zap.L().Error("Failed to find user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	zap.L().Info("User found",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, response)
}

func (h *userHandler) DeleteUser(context *gin.Context) {
	userID, _ := context.Get("userID")

	var user entities.User

	result := common.DB.First(&user, userID)
	if result.Error != nil {
		zap.L().Error("Failed to find user to delete",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	result = common.DB.Delete(&user)
	if result.Error != nil {
		zap.L().Error("Failed to delete user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	zap.L().Info("User deleted",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, gin.H{"message": "You successfully deleted your account."})
}

func (h *userHandler) ChangePassword(context *gin.Context) {
	userID, _ := context.Get("userID")

	var user entities.User

	result := common.DB.First(&user, userID)
	if result.Error != nil {
		zap.L().Error("Failed to find user to change password",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var request models.ChangePasswordRequest
	err := context.ShouldBindJSON(&request)
	if err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(request); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !common.CheckPasswordHash(request.OldPassword, user.Password) {
		zap.L().Error("Invalid password",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	hash, err := common.HashPassword(request.NewPassword)
	if err != nil {
		zap.L().Error("Failed to hash password",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hash

	result = common.DB.Save(&user)
	if result.Error != nil {
		zap.L().Error("Failed to save user",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	zap.L().Info("Password changed successfully",
		zap.Uint64("user ID", user.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, gin.H{"message": "You successfully changed your password."})
}
