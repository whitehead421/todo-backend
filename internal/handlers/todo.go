package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
	"go.uber.org/zap"
)

var validate *validator.Validate

func CreateTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	var todoRequest models.TodoRequest

	if err := context.ShouldBindJSON(&todoRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.Struct(todoRequest); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := entities.Todo{
		Description: todoRequest.Description,
		Status:      string(models.Pending),
		UserID:      userID.(uint64),
	}

	result := common.DB.Create(&todo)
	if result.Error != nil {
		zap.L().Error("Failed to create todo",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := models.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	zap.L().Info("Todo created successfully",
		zap.Uint64("todo ID", todo.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, todoResponse)
}

func ReadTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")
	var todo entities.Todo

	result := common.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo)
	if result.Error != nil {
		zap.L().Error("Failed to find todo",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := models.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	zap.L().Info("Todo found successfully",
		zap.Uint64("todo ID", todo.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, todoResponse)
}

func UpdateTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")
	// Check if ID is valid
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zap.L().Error("Invalid ID",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if todo to update exists, if not return 404
	result := common.DB.Where("id = ? AND user_id = ?", ID, userID).First(&entities.Todo{})
	if result.Error != nil {
		zap.L().Error("Failed to find todo to update",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	var todoUpdateRequest models.TodoUpdateRequest

	if err := context.ShouldBindJSON(&todoUpdateRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.StructPartial(todoUpdateRequest); err != nil {
		zap.L().Error("Validation error",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := entities.Todo{
		ID:          ID,
		Description: todoUpdateRequest.Description,
		Status:      string(todoUpdateRequest.Status),
	}

	// Update todo
	result = common.DB.Save(&todo)
	if result.Error != nil {
		zap.L().Error("Failed to update todo",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := models.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	zap.L().Info("Todo updated successfully",
		zap.Uint64("todo ID", todo.ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, todoResponse)
}

func DeleteTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")

	// Check if ID is valid
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zap.L().Error("Invalid ID",
			zap.Error(err),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if todo to delete exists, if not return 404
	result := common.DB.Where("id = ? AND user_id = ?", ID, userID).First(&entities.Todo{})
	if result.Error != nil {
		zap.L().Error("Failed to find todo to delete",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	result = common.DB.Delete(&entities.Todo{ID: ID})
	if result.Error != nil {
		zap.L().Error("Failed to delete todo",
			zap.Error(result.Error),
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	zap.L().Info("Todo deleted successfully",
		zap.Uint64("todo ID", ID),
		zap.String("url path", context.Request.URL.Path),
	)

	context.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
