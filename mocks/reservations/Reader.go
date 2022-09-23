// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
	reservations "github.ibm.com/skol/atkcli/pkg/reservations"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

// Read provides a mock function with given fields: _a0
func (_m *Reader) Read(_a0 io.Reader) (reservations.TZReservation, error) {
	ret := _m.Called(_a0)

	var r0 reservations.TZReservation
	if rf, ok := ret.Get(0).(func(io.Reader) reservations.TZReservation); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(reservations.TZReservation)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadAll provides a mock function with given fields: _a0
func (_m *Reader) ReadAll(_a0 io.Reader) ([]reservations.TZReservation, error) {
	ret := _m.Called(_a0)

	var r0 []reservations.TZReservation
	if rf, ok := ret.Get(0).(func(io.Reader) []reservations.TZReservation); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reservations.TZReservation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReader interface {
	mock.TestingT
	Cleanup(func())
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReader(t mockConstructorTestingTNewReader) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
