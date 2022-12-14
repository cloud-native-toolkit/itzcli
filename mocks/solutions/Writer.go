// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
	solutions "github.com/cloud-native-toolkit/itzcli/pkg/solutions"
)

// Writer is an autogenerated mock type for the Writer type
type Writer struct {
	mock.Mock
}

// Write provides a mock function with given fields: _a0, _a1
func (_m *Writer) Write(_a0 io.Writer, _a1 solutions.Solution) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, solutions.Solution) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteAll provides a mock function with given fields: _a0, _a1
func (_m *Writer) WriteAll(_a0 io.Writer, _a1 []solutions.Solution) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, []solutions.Solution) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteFilter provides a mock function with given fields: _a0, _a1, _a2
func (_m *Writer) WriteFilter(_a0 io.Writer, _a1 []solutions.Solution, _a2 solutions.Filter) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, []solutions.Solution, solutions.Filter) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewWriter interface {
	mock.TestingT
	Cleanup(func())
}

// NewWriter creates a new instance of Writer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWriter(t mockConstructorTestingTNewWriter) *Writer {
	mock := &Writer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
