package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/mocks"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
)

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uint64(1)

	todoHandler := handlers.NewTodoHandler()

	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		c.Set("userID", userID) // Set userID in context
		todoHandler.CreateTodo(c)
	})

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(mockTodoHandler *mocks.TodoHandler)
		expectedStatus int
	}{
		{
			name:           "Successfull Todo Creation",
			requestBody:    models.TodoRequest{Description: "Test Todo"},
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Validation Error",
			requestBody:    models.TodoRequest{Description: ""},
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockTodoHandler := new(mocks.TodoHandler)
			tt.mockBehavior(mockTodoHandler)

			// Create request body
			var reqBody []byte
			var err error
			if body, ok := tt.requestBody.(models.TodoRequest); ok {
				reqBody, err = json.Marshal(body)
				assert.NoError(t, err)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTodoHandler.AssertExpectations(t)
		})
	}
}

func TestReadTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uint64(1)

	todoHandler := handlers.NewTodoHandler()

	router := gin.Default()
	router.GET("/:id", func(c *gin.Context) {
		c.Set("userID", userID) // Set userID in context
		todoHandler.ReadTodo(c)
	})

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		id             uint64
		mockBehavior   func(mockTodoHandler *mocks.TodoHandler)
		expectedStatus int
	}{
		{
			name: "Successfull Todo Read",
			id:   1,
			mockBehavior: func(mockTodoHandler *mocks.TodoHandler) {
				common.DB.Create(&entities.Todo{ID: 1, Description: "Test Todo", Status: "pending", UserID: 1})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Todo Not Found",
			id:             4,
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockTodoHandler := new(mocks.TodoHandler)
			tt.mockBehavior(mockTodoHandler)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", tt.id), nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTodoHandler.AssertExpectations(t)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uint64(1)

	todoHandler := handlers.NewTodoHandler()

	router := gin.Default()
	router.PUT("/:id", func(c *gin.Context) {
		c.Set("userID", userID) // Set userID in context
		todoHandler.UpdateTodo(c)
	})

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		id             uint64
		requestBody    interface{}
		mockBehavior   func(mockTodoHandler *mocks.TodoHandler)
		expectedStatus int
	}{
		{
			name:        "Successfull Todo Update",
			id:          1,
			requestBody: models.TodoUpdateRequest{Description: "Test Todo", Status: "pending"},
			mockBehavior: func(mockTodoHandler *mocks.TodoHandler) {
				common.DB.Create(&entities.User{ID: 1, Email: "test@test.com", Name: "Test User"})
				common.DB.Create(&entities.Todo{ID: 1, Description: "Test Todo", Status: "pending", UserID: 1})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			id:             1,
			requestBody:    "invalid json",
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Todo Not Found",
			id:             4,
			requestBody:    models.TodoUpdateRequest{Description: "Test Todo", Status: "pending"},
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockTodoHandler := new(mocks.TodoHandler)
			tt.mockBehavior(mockTodoHandler)

			// Create request body
			var reqBody []byte
			var err error
			if body, ok := tt.requestBody.(models.TodoUpdateRequest); ok {
				reqBody, err = json.Marshal(body)
				assert.NoError(t, err)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/%d", tt.id), bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTodoHandler.AssertExpectations(t)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uint64(1)

	todoHandler := handlers.NewTodoHandler()

	router := gin.Default()
	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("userID", userID) // Set userID in context
		todoHandler.DeleteTodo(c)
	})

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		id             uint64
		mockBehavior   func(mockTodoHandler *mocks.TodoHandler)
		expectedStatus int
	}{
		{
			name: "Successfull Todo Deletion",
			id:   1,
			mockBehavior: func(mockTodoHandler *mocks.TodoHandler) {
				common.DB.Create(&entities.Todo{ID: 1, Description: "Test Todo", Status: "pending", UserID: 1})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Todo Not Found",
			id:             4,
			mockBehavior:   func(mockTodoHandler *mocks.TodoHandler) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockTodoHandler := new(mocks.TodoHandler)
			tt.mockBehavior(mockTodoHandler)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/%d", tt.id), nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTodoHandler.AssertExpectations(t)
		})
	}
}
