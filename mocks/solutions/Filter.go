// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	solutions "github.com/cloud-native-toolkit/itzcli/pkg/solutions"
)

// Filter is an autogenerated mock type for the Filter type
type Filter struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Filter) Execute(_a0 solutions.Solution) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(solutions.Solution) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewFilter interface {
	mock.TestingT
	Cleanup(func())
}

// NewFilter creates a new instance of Filter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilter(t mockConstructorTestingTNewFilter) *Filter {
	mock := &Filter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
