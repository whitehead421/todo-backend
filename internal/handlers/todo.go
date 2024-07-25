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
	"gorm.io/gorm"
)

type TodoHandler interface {
	CreateTodo(context *gin.Context)
	ReadTodo(context *gin.Context)
	UpdateTodo(context *gin.Context)
	DeleteTodo(context *gin.Context)
}

type todoHandler struct {
	validate *validator.Validate
}

func NewTodoHandler() TodoHandler {
	return &todoHandler{
		validate: validator.New(),
	}
}

func (h *todoHandler) CreateTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	var todoRequest models.TodoRequest

	if err := context.ShouldBindJSON(&todoRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(todoRequest); err != nil {
		zap.L().Error("Validation error",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
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
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
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

func (h *todoHandler) ReadTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")
	var todo entities.Todo

	result := common.DB.First(&todo, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			zap.L().Error("Todo not found",
				zap.String("url path", context.Request.URL.Path),
				zap.Error(result.Error),
			)
			context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		zap.L().Error("Failed to find todo",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if todo.UserID != userID {
		zap.L().Error("User does not have permission to access this todo",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this todo"})
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

func (h *todoHandler) UpdateTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")
	// Check if ID is valid
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zap.L().Error("Invalid ID",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo entities.Todo

	// Check if todo to update exists, if not return 404
	result := common.DB.First(&todo, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			zap.L().Error("Todo not found",
				zap.String("url path", context.Request.URL.Path),
				zap.Error(result.Error),
			)
			context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		zap.L().Error("Failed to find todo to update",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
		)

		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if todo.UserID != userID {
		zap.L().Error("User does not have permission to access this todo",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this todo"})
		return
	}

	var todoUpdateRequest models.TodoUpdateRequest

	if err := context.ShouldBindJSON(&todoUpdateRequest); err != nil {
		zap.L().Error("Failed to bind JSON",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.StructPartial(todoUpdateRequest); err != nil {
		zap.L().Error("Validation error",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo = entities.Todo{
		ID:          ID,
		Description: todoUpdateRequest.Description,
		Status:      string(todoUpdateRequest.Status),
		UserID:      userID.(uint64),
	}

	// Update todo
	result = common.DB.Save(&todo)
	if result.Error != nil {
		zap.L().Error("Failed to update todo",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
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

func (h *todoHandler) DeleteTodo(context *gin.Context) {
	// Ignoring exists check as we are using authentication middleware, so it should always exist
	userID, _ := context.Get("userID")

	id := context.Param("id")

	// Check if ID is valid
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zap.L().Error("Invalid ID",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(err),
		)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo entities.Todo

	// Check if todo to delete exists, if not return 404
	result := common.DB.First(&todo, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			zap.L().Error("Todo not found",
				zap.String("url path", context.Request.URL.Path),
				zap.Error(result.Error),
			)
			context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		zap.L().Error("Failed to find todo",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
		)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if todo.UserID != userID {
		zap.L().Error("User does not have permission to access this todo",
			zap.String("url path", context.Request.URL.Path),
		)
		context.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this todo"})
		return
	}

	result = common.DB.Delete(&entities.Todo{ID: ID})
	if result.Error != nil {
		zap.L().Error("Failed to delete todo",
			zap.String("url path", context.Request.URL.Path),
			zap.Error(result.Error),
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
