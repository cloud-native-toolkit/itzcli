// Code generated by mockery v2.25.1. DO NOT EDIT.

package mocks

import (
	pkg "github.com/cloud-native-toolkit/itzcli/pkg"
	mock "github.com/stretchr/testify/mock"
)

// PipelineServiceClient is an autogenerated mock type for the PipelineServiceClient type
type PipelineServiceClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *PipelineServiceClient) Get(id string) (*pkg.Pipeline, error) {
	ret := _m.Called(id)

	var r0 *pkg.Pipeline
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*pkg.Pipeline, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *pkg.Pipeline); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.Pipeline)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *PipelineServiceClient) GetAll() ([]*pkg.Pipeline, error) {
	ret := _m.Called()

	var r0 []*pkg.Pipeline
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*pkg.Pipeline, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*pkg.Pipeline); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pkg.Pipeline)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPipelineServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewPipelineServiceClient creates a new instance of PipelineServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPipelineServiceClient(t mockConstructorTestingTNewPipelineServiceClient) *PipelineServiceClient {
	mock := &PipelineServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}