// Code generated by mockery v2.25.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MarshallerFunc is an autogenerated mock type for the MarshallerFunc type
type MarshallerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: b
func (_m *MarshallerFunc) Execute(b []byte) (interface{}, error) {
	ret := _m.Called(b)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (interface{}, error)); ok {
		return rf(b)
	}
	if rf, ok := ret.Get(0).(func([]byte) interface{}); ok {
		r0 = rf(b)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMarshallerFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewMarshallerFunc creates a new instance of MarshallerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMarshallerFunc(t mockConstructorTestingTNewMarshallerFunc) *MarshallerFunc {
	mock := &MarshallerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
