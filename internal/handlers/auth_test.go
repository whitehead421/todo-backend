package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/mocks"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/models"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := handlers.NewAuthHandler()

	router := gin.Default()
	router.POST("/register", authHandler.Register)

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(mockAuthHandler *mocks.AuthHandler)
		expectedStatus int
	}{
		{
			name:           "Failed to bind JSON",
			requestBody:    "invalid json",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: models.RegisterRequest{
				Email:    "invalid-email",
				Name:     "",
				Password: "",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Successful registration",
			requestBody: models.RegisterRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password",
				Confirm:  "password",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockAuthHandler := new(mocks.AuthHandler)
			tt.mockBehavior(mockAuthHandler)

			// Create request body
			var reqBody []byte
			var err error
			if body, ok := tt.requestBody.(models.RegisterRequest); ok {
				reqBody, err = json.Marshal(body)
				assert.NoError(t, err)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthHandler.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := handlers.NewAuthHandler()

	router := gin.Default()
	router.POST("/login", authHandler.Login)

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(mockAuthHandler *mocks.AuthHandler)
		expectedStatus int
	}{
		{
			name:           "Failed to bind JSON",
			requestBody:    "invalid json",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: models.LoginRequest{
				Email:    "invalid-email",
				Password: "",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failed to find user",
			requestBody: models.LoginRequest{
				Email:    "unknown_user@test.com",
				Password: "password",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Email:    "test@test.com",
				Password: "password123",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockAuthHandler := new(mocks.AuthHandler)
			tt.mockBehavior(mockAuthHandler)

			// Create request body
			var reqBody []byte
			var err error
			if body, ok := tt.requestBody.(models.LoginRequest); ok {
				reqBody, err = json.Marshal(body)
				assert.NoError(t, err)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthHandler.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := handlers.NewAuthHandler()

	router := gin.Default()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("userID", uint64(1))
		authHandler.Logout(c)
	})

	// Mock Redis setup
	db, mock := redismock.NewClientMock()
	common.SetRedisClient(db)

	tests := []struct {
		name           string
		mockBehavior   func(mock redismock.ClientMock)
		expectedStatus int
	}{
		{
			name: "Successful logout",
			mockBehavior: func(mock redismock.ClientMock) {
				mock.ExpectDel("token").SetVal(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failed to delete token",
			mockBehavior: func(mock redismock.ClientMock) {
				mock.ExpectDel("token").SetErr(nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock Redis client
			tt.mockBehavior(mock)

			req, err := http.NewRequest(http.MethodPost, "/logout", nil)
			assert.NoError(t, err)

			// Add Authorization header with Bearer token
			req.Header.Set("Authorization", "Bearer token")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}
