package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/whitehead421/todo-backend/internal/handlers"
	"github.com/whitehead421/todo-backend/mocks"
	"github.com/whitehead421/todo-backend/pkg/common"
	"github.com/whitehead421/todo-backend/pkg/entities"
	"github.com/whitehead421/todo-backend/pkg/models"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKafkaWriter := &mocks.KafkaWriter{}
	authHandler := handlers.NewAuthHandler(mockKafkaWriter)

	router := gin.Default()
	router.POST("/register", authHandler.Register)

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(mockAuthHandler *mocks.AuthHandler, mockKafkaWriter *mocks.KafkaWriter)
		expectedStatus int
	}{
		{
			name:           "Failed to bind JSON",
			requestBody:    "invalid json",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler, mockKafkaWriter *mocks.KafkaWriter) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: models.RegisterRequest{
				Email:    "invalid-email",
				Name:     "",
				Password: "",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler, mockKafkaWriter *mocks.KafkaWriter) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Email already exists",
			requestBody: models.RegisterRequest{
				Email:    "registered@example.com",
				Name:     "Test User",
				Password: "password",
				Confirm:  "password",
			},
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler, mockKafkaWriter *mocks.KafkaWriter) {
				common.DB.Create(&entities.User{
					Email:    "registered@example.com",
					Name:     "Test User",
					Password: "password",
				})
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Successful registration",
			requestBody: models.RegisterRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password",
				Confirm:  "password",
			},
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler, mockKafkaWriter *mocks.KafkaWriter) {
				mockKafkaWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockAuthHandler := new(mocks.AuthHandler)
			tt.mockBehavior(mockAuthHandler, mockKafkaWriter)

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

	mockKafkaWriter := &mocks.KafkaWriter{}
	authHandler := handlers.NewAuthHandler(mockKafkaWriter)

	router := gin.Default()
	router.POST("/login", authHandler.Login)

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	// Mock Redis setup
	db, mockRedis := redismock.NewClientMock()
	common.SetRedisClient(db)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock)
		expectedStatus int
	}{
		{
			name:           "Failed to bind JSON",
			requestBody:    "invalid json",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: models.LoginRequest{
				Email:    "invalid-email",
				Password: "",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failed to find user",
			requestBody: models.LoginRequest{
				Email:    "unknown_user@example.com",
				Password: "password",
			},
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Account is not verified",
			requestBody: models.LoginRequest{
				Email:    "unverified@example.com",
				Password: "password",
			},
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {
				common.DB.Create(&entities.User{
					Email:    "unverified@example.com",
					Name:     "Unverified User",
					Password: common.HashPassword("password"),
					Verified: false,
				})
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Email:    "invalidpassword@example.com",
				Password: "invalid-password",
			},
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {
				common.DB.Create(&entities.User{
					Email:    "invalidpassword@example.com",
					Name:     "Invalid Password",
					Password: common.HashPassword("password"),
					Verified: true,
				})
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failed to set token to Redis",
			requestBody: models.LoginRequest{
				Email:    "failedredis@example.com",
				Password: "password",
			},
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler, mockRedis redismock.ClientMock) {
				common.DB.Create(&entities.User{
					Email:    "failedredis@example.com",
					Name:     "Failed Redis",
					Password: common.HashPassword("password"),
					Verified: true,
				})
				mockRedis.ExpectSet("failedredis@example.com", "token", time.Hour).SetErr(errors.New("failed to set token"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockAuthHandler := new(mocks.AuthHandler)
			tt.mockBehavior(mockAuthHandler, mockRedis)

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

	mockKafkaWriter := &mocks.KafkaWriter{}
	authHandler := handlers.NewAuthHandler(mockKafkaWriter)

	router := gin.Default()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("userID", uint64(1))
		authHandler.Logout(c)
	})

	// Mock Redis setup
	db, mockRedis := redismock.NewClientMock()
	common.SetRedisClient(db)

	tests := []struct {
		name           string
		mockBehavior   func(mockRedis redismock.ClientMock)
		expectedStatus int
		Header         map[string]string
	}{
		{
			name: "Successful logout",
			mockBehavior: func(mockRedis redismock.ClientMock) {
				mockRedis.ExpectDel("token").SetVal(1)
			},
			expectedStatus: http.StatusOK,
			Header: map[string]string{
				"Authorization": "Bearer token",
			},
		},
		{
			name: "Failed to delete token",
			mockBehavior: func(mockRedis redismock.ClientMock) {
				mockRedis.ExpectDel("token").SetErr(nil)
			},
			expectedStatus: http.StatusInternalServerError,
			Header: map[string]string{
				"Authorization": "Bearer token",
			},
		},
		{
			name:           "Authorization header is missing",
			mockBehavior:   func(mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusUnauthorized,
			Header:         map[string]string{},
		},
		{
			name:           "Authorization header format is invalid",
			mockBehavior:   func(mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusUnauthorized,
			Header: map[string]string{
				"Authorization": "invalid-format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock Redis client
			tt.mockBehavior(mockRedis)

			req, err := http.NewRequest(http.MethodPost, "/logout", nil)
			assert.NoError(t, err)

			// Add headers to request
			for key, value := range tt.Header {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Ensure all expectations were met
			if err := mockRedis.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKafkaWriter := &mocks.KafkaWriter{}
	authHandler := handlers.NewAuthHandler(mockKafkaWriter)

	router := gin.Default()
	router.GET("/verify", authHandler.Verify)

	// Setup test database
	testDB := common.SetupTestDB()
	common.SetDB(testDB) // Set the mock database for testing

	uuidToken := common.GenerateUUID()

	tests := []struct {
		name           string
		mockBehavior   func(mockAuthHandler *mocks.AuthHandler)
		expectedStatus int
		queryParams    map[string]string
	}{
		{
			name: "Successful verification",
			mockBehavior: func(mockAuthHandler *mocks.AuthHandler) {
				common.DB.Create(&entities.User{
					Email:       "test@example.com",
					Name:        "Test User",
					Password:    common.HashPassword("password"),
					Verified:    false,
					VerifyToken: uuidToken,
				})
			},
			expectedStatus: http.StatusOK,
			queryParams: map[string]string{
				"token": uuidToken,
			},
		},
		{
			name:           "Verification token not found",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusBadRequest,
			queryParams:    map[string]string{},
		},
		{
			name:           "Failed to find user",
			mockBehavior:   func(mockAuthHandler *mocks.AuthHandler) {},
			expectedStatus: http.StatusNotFound,
			queryParams: map[string]string{
				"token": "invalid-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock handler
			mockAuthHandler := new(mocks.AuthHandler)
			tt.mockBehavior(mockAuthHandler)

			req, err := http.NewRequest(http.MethodGet, "/verify", nil)
			assert.NoError(t, err)

			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthHandler.AssertExpectations(t)
		})
	}
}

func TestAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKafkaWriter := &mocks.KafkaWriter{}
	authHandler := handlers.NewAuthHandler(mockKafkaWriter)

	router := gin.Default()
	router.POST("/authorize", authHandler.Authorize)

	// Mock Redis setup
	db, mockRedis := redismock.NewClientMock()
	common.SetRedisClient(db)

	tests := []struct {
		name           string
		mockBehavior   func(mockRedis redismock.ClientMock)
		expectedStatus int
		Header         map[string]string
	}{
		{
			name:           "Authorization header is missing",
			mockBehavior:   func(mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusUnauthorized,
			Header:         map[string]string{},
		},
		{
			name:           "Authorization header format is invalid",
			mockBehavior:   func(mockRedis redismock.ClientMock) {},
			expectedStatus: http.StatusUnauthorized,
			Header: map[string]string{
				"Authorization": "invalid-format",
			},
		},
		{
			name: "Failed to get token from Redis",
			mockBehavior: func(mockRedis redismock.ClientMock) {
				mockRedis.ExpectGet("token").SetErr(errors.New("failed to get token"))
			},
			expectedStatus: http.StatusUnauthorized,
			Header: map[string]string{
				"Authorization": "Bearer token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock Redis client
			tt.mockBehavior(mockRedis)

			req, err := http.NewRequest(http.MethodPost, "/authorize", nil)
			assert.NoError(t, err)

			// Add headers to request
			for key, value := range tt.Header {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
