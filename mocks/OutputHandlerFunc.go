// Code generated by mockery v2.25.1. DO NOT EDIT.

package mocks

import (
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
)

// OutputHandlerFunc is an autogenerated mock type for the OutputHandlerFunc type
type OutputHandlerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: out
func (_m *OutputHandlerFunc) Execute(out *bytes.Buffer) error {
	ret := _m.Called(out)

	var r0 error
	if rf, ok := ret.Get(0).(func(*bytes.Buffer) error); ok {
		r0 = rf(out)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOutputHandlerFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewOutputHandlerFunc creates a new instance of OutputHandlerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOutputHandlerFunc(t mockConstructorTestingTNewOutputHandlerFunc) *OutputHandlerFunc {
	mock := &OutputHandlerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
