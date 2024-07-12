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
)

var validate *validator.Validate

// @Summary Create todo
// @Description Create todo
// @Produce  json
// @Param todo body entities.TodoRequest true "Todo Request"
// @Success 200 {object} entities.TodoResponse
// @Router / [post]
func CreateTodo(context *gin.Context) {
	var todoRequest entities.TodoRequest

	if err := context.ShouldBindJSON(&todoRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.Struct(todoRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := models.Todo{
		Description: todoRequest.Description,
		Status:      string(entities.Pending),
	}

	result := common.DB.Create(&todo)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := entities.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	context.JSON(http.StatusOK, todoResponse)
}

// @Summary Get todo
// @Description Get todo
// @Produce  json
// @Param id path string true "Todo ID"
// @Success 200 {object} entities.TodoResponse
// @Router /{id} [get]
func ReadTodo(context *gin.Context) {
	id := context.Param("id")
	var todo models.Todo

	result := common.DB.First(&todo, id)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := entities.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	context.JSON(http.StatusOK, todoResponse)
}

// @Summary Update todo
// @Description Update todo
// @Produce  json
// @Param id path string true "Todo ID"
// @Param todo body entities.TodoUpdateRequest true "Todo Request"
// @Success 200 {object} entities.TodoResponse
// @Router /{id} [put]
func UpdateTodo(context *gin.Context) {
	id := context.Param("id")
	// Check if ID is valid
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if todo to update exists, if not return 404
	result := common.DB.First(&models.Todo{ID: ID})
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	var todoUpdateRequest entities.TodoUpdateRequest

	if err := context.ShouldBindJSON(&todoUpdateRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate = validator.New()
	if err := validate.StructPartial(todoUpdateRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := models.Todo{
		ID:          ID,
		Description: todoUpdateRequest.Description,
		Status:      string(todoUpdateRequest.Status),
	}

	// Update todo
	result = common.DB.Save(&todo)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	todoResponse := entities.TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	context.JSON(http.StatusOK, todoResponse)
}

// @Summary Delete todo
// @Description Delete todo
// @Produce  json
// @Router / [delete]
func DeleteTodo(context *gin.Context) {
	context.JSON(http.StatusOK, "Delete Todo")
}
