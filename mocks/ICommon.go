// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// ICommon is an autogenerated mock type for the ICommon type
type ICommon struct {
	mock.Mock
}

// BlacklistToken provides a mock function with given fields: tokenString, expiration, ctx
func (_m *ICommon) BlacklistToken(tokenString string, expiration time.Duration, ctx context.Context) error {
	ret := _m.Called(tokenString, expiration, ctx)

	if len(ret) == 0 {
		panic("no return value specified for BlacklistToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Duration, context.Context) error); ok {
		r0 = rf(tokenString, expiration, ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckPasswordHash provides a mock function with given fields: password, hash
func (_m *ICommon) CheckPasswordHash(password string, hash string) bool {
	ret := _m.Called(password, hash)

	if len(ret) == 0 {
		panic("no return value specified for CheckPasswordHash")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(password, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CreateToken provides a mock function with given fields: id
func (_m *ICommon) CreateToken(id uint64) (string, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for CreateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (string, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HashPassword provides a mock function with given fields: password
func (_m *ICommon) HashPassword(password string) string {
	ret := _m.Called(password)

	if len(ret) == 0 {
		panic("no return value specified for HashPassword")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsTokenBlacklisted provides a mock function with given fields: tokenString, ctx
func (_m *ICommon) IsTokenBlacklisted(tokenString string, ctx context.Context) (bool, error) {
	ret := _m.Called(tokenString, ctx)

	if len(ret) == 0 {
		panic("no return value specified for IsTokenBlacklisted")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, context.Context) (bool, error)); ok {
		return rf(tokenString, ctx)
	}
	if rf, ok := ret.Get(0).(func(string, context.Context) bool); ok {
		r0 = rf(tokenString, ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, context.Context) error); ok {
		r1 = rf(tokenString, ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICommon creates a new instance of ICommon. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICommon(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICommon {
	mock := &ICommon{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
