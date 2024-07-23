// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// UserHandler is an autogenerated mock type for the UserHandler type
type UserHandler struct {
	mock.Mock
}

// ChangePassword provides a mock function with given fields: context
func (_m *UserHandler) ChangePassword(context *gin.Context) {
	_m.Called(context)
}

// DeleteUser provides a mock function with given fields: context
func (_m *UserHandler) DeleteUser(context *gin.Context) {
	_m.Called(context)
}

// GetUser provides a mock function with given fields: context
func (_m *UserHandler) GetUser(context *gin.Context) {
	_m.Called(context)
}

// NewUserHandler creates a new instance of UserHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserHandler {
	mock := &UserHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
