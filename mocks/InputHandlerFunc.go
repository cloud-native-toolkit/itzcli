// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
)

// InputHandlerFunc is an autogenerated mock type for the InputHandlerFunc type
type InputHandlerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: in
func (_m *InputHandlerFunc) Execute(in *bytes.Buffer) error {
	ret := _m.Called(in)

	var r0 error
	if rf, ok := ret.Get(0).(func(*bytes.Buffer) error); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewInputHandlerFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewInputHandlerFunc creates a new instance of InputHandlerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInputHandlerFunc(t mockConstructorTestingTNewInputHandlerFunc) *InputHandlerFunc {
	mock := &InputHandlerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
