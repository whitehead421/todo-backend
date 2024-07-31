// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	kafka "github.com/segmentio/kafka-go"

	mock "github.com/stretchr/testify/mock"
)

// KafkaWriter is an autogenerated mock type for the KafkaWriter type
type KafkaWriter struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *KafkaWriter) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteMessages provides a mock function with given fields: ctx, msgs
func (_m *KafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	_va := make([]interface{}, len(msgs))
	for _i := range msgs {
		_va[_i] = msgs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for WriteMessages")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...kafka.Message) error); ok {
		r0 = rf(ctx, msgs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewKafkaWriter creates a new instance of KafkaWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKafkaWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *KafkaWriter {
	mock := &KafkaWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
