// Code generated by mockery v2.25.1. DO NOT EDIT.

package mocks

import (
	io "io"

	reservations "github.com/cloud-native-toolkit/itzcli/pkg/reservations"
	mock "github.com/stretchr/testify/mock"
)

// ReservationWriter is an autogenerated mock type for the ReservationWriter type
type ReservationWriter struct {
	mock.Mock
}

// WriteMany provides a mock function with given fields: w, rezs
func (_m *ReservationWriter) WriteMany(w io.Writer, rezs []reservations.TZReservation) error {
	ret := _m.Called(w, rezs)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, []reservations.TZReservation) error); ok {
		r0 = rf(w, rezs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteOne provides a mock function with given fields: w, rez
func (_m *ReservationWriter) WriteOne(w io.Writer, rez reservations.TZReservation) error {
	ret := _m.Called(w, rez)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, reservations.TZReservation) error); ok {
		r0 = rf(w, rez)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewReservationWriter interface {
	mock.TestingT
	Cleanup(func())
}

// NewReservationWriter creates a new instance of ReservationWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReservationWriter(t mockConstructorTestingTNewReservationWriter) *ReservationWriter {
	mock := &ReservationWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
