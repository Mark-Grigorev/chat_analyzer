// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// LLMClient is an autogenerated mock type for the LLMClient type
type LLMClient struct {
	mock.Mock
}

// GetLLMResponseAboutMsg provides a mock function with given fields: ctx, promt
func (_m *LLMClient) GetLLMResponseAboutMsg(ctx context.Context, promt string) (string, error) {
	ret := _m.Called(ctx, promt)

	if len(ret) == 0 {
		panic("no return value specified for GetLLMResponseAboutMsg")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, promt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, promt)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, promt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLLMClient creates a new instance of LLMClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLLMClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *LLMClient {
	mock := &LLMClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
