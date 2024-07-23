package handlers_test

import (
	"bytes"
	"encoding/json"
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

func TestGetUser(t *testing.T) {
	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		userID         interface{}
		mockBehavior   func(mockUserHandler *mocks.UserHandler)
		expectedStatus int
	}{
		{
			name:   "Successfull User Retrieval",
			userID: uint64(1),
			mockBehavior: func(mockUserHandler *mocks.UserHandler) {
				common.DB.Create(&entities.User{
					ID:    uint64(1),
					Name:  "Test User",
					Email: "test@test.com",
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User Not Found",
			userID:         uint64(999),
			mockBehavior:   func(mockUserHandler *mocks.UserHandler) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router for every test case
			gin.SetMode(gin.TestMode)

			userHandler := handlers.NewUserHandler()

			// Setup router
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				c.Next()
			})
			router.GET("/", userHandler.GetUser)

			// Setup mock handler
			mockUserHandler := new(mocks.UserHandler)
			tt.mockBehavior(mockUserHandler)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			// Create a new context with the specified userID
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserHandler.AssertExpectations(t)

			router = nil
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		userID         interface{}
		mockBehavior   func(mockUserHandler *mocks.UserHandler)
		expectedStatus int
	}{
		{
			name:   "Successfull User Deletion",
			userID: uint64(1),
			mockBehavior: func(mockUserHandler *mocks.UserHandler) {
				common.DB.Create(&entities.User{
					ID:    uint64(1),
					Name:  "Test User",
					Email: "test@test.com",
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User Not Found",
			userID:         uint64(999),
			mockBehavior:   func(mockUserHandler *mocks.UserHandler) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router for every test case
			gin.SetMode(gin.TestMode)

			userHandler := handlers.NewUserHandler()

			// Setup router
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				c.Next()
			})
			router.DELETE("/", userHandler.DeleteUser)

			// Setup mock handler
			mockUserHandler := new(mocks.UserHandler)
			tt.mockBehavior(mockUserHandler)

			req, err := http.NewRequest(http.MethodDelete, "/", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			// Create a new context with the specified userID
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserHandler.AssertExpectations(t)

			router = nil
		})
	}
}

func TestChangePassword(t *testing.T) {
	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		userID         interface{}
		requestBody    interface{}
		mockBehavior   func()
		expectedStatus int
	}{
		{
			name:   "Successful Password Change",
			userID: uint64(1),
			requestBody: models.ChangePasswordRequest{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockBehavior: func() {
				// Add a mock user to the database
				hashedOldPassword := common.HashPassword("oldpassword")
				common.DB.Create(&entities.User{
					ID:       1,
					Email:    "test@example.com",
					Name:     "Test User",
					Password: hashedOldPassword,
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid Old Password",
			userID: uint64(1),
			requestBody: models.ChangePasswordRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword",
			},
			mockBehavior: func() {
				// Add a mock user with a different password
				hashedOldPassword := common.HashPassword("oldpassword")
				common.DB.Create(&entities.User{
					ID:       1,
					Email:    "test@example.com",
					Name:     "Test User",
					Password: hashedOldPassword,
				})
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid JSON",
			userID:         uint64(1),
			requestBody:    "invalid json",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "User Not Found",
			userID: uint64(999),
			requestBody: models.ChangePasswordRequest{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router for every test case
			gin.SetMode(gin.TestMode)

			userHandler := handlers.NewUserHandler()

			router := gin.Default()
			router.Use(func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				c.Next()
			})
			router.PUT("/", userHandler.ChangePassword)

			// Apply mock behavior
			tt.mockBehavior()

			// Create request body
			var reqBody []byte
			var err error
			if body, ok := tt.requestBody.(models.ChangePasswordRequest); ok {
				reqBody, err = json.Marshal(body)
				assert.NoError(t, err)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
