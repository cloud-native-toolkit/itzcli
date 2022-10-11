// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ParamBuilderFunc is an autogenerated mock type for the ParamBuilderFunc type
type ParamBuilderFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields:
func (_m *ParamBuilderFunc) Execute() map[string]string {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

type mockConstructorTestingTNewParamBuilderFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewParamBuilderFunc creates a new instance of ParamBuilderFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewParamBuilderFunc(t mockConstructorTestingTNewParamBuilderFunc) *ParamBuilderFunc {
	mock := &ParamBuilderFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}