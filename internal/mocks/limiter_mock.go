// Code generated by MockGen. DO NOT EDIT.
// Source: limiter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLimiter is a mock of Limiter interface
type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
}

// MockLimiterMockRecorder is the mock recorder for MockLimiter
type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

// NewMockLimiter creates a new mock instance
func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return m.recorder
}

// CheckLimits mocks base method
func (m *MockLimiter) CheckLimits(ctx context.Context, login, password, ip string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLimits", ctx, login, password, ip)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLimits indicates an expected call of CheckLimits
func (mr *MockLimiterMockRecorder) CheckLimits(ctx, login, password, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLimits", reflect.TypeOf((*MockLimiter)(nil).CheckLimits), ctx, login, password, ip)
}

// DropBuckets mocks base method
func (m *MockLimiter) DropBuckets(ctx context.Context, login, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropBuckets", ctx, login, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// DropBuckets indicates an expected call of DropBuckets
func (mr *MockLimiterMockRecorder) DropBuckets(ctx, login, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropBuckets", reflect.TypeOf((*MockLimiter)(nil).DropBuckets), ctx, login, ip)
}
