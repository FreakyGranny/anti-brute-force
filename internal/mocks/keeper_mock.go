// Code generated by MockGen. DO NOT EDIT.
// Source: keeper.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	net "net"
	reflect "reflect"
	time "time"
)

// MockIPKeeper is a mock of IPKeeper interface
type MockIPKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockIPKeeperMockRecorder
}

// MockIPKeeperMockRecorder is the mock recorder for MockIPKeeper
type MockIPKeeperMockRecorder struct {
	mock *MockIPKeeper
}

// NewMockIPKeeper creates a new mock instance
func NewMockIPKeeper(ctrl *gomock.Controller) *MockIPKeeper {
	mock := &MockIPKeeper{ctrl: ctrl}
	mock.recorder = &MockIPKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIPKeeper) EXPECT() *MockIPKeeperMockRecorder {
	return m.recorder
}

// GetBlacklist mocks base method
func (m *MockIPKeeper) GetBlacklist() []*net.IPNet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlacklist")
	ret0, _ := ret[0].([]*net.IPNet)
	return ret0
}

// GetBlacklist indicates an expected call of GetBlacklist
func (mr *MockIPKeeperMockRecorder) GetBlacklist() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlacklist", reflect.TypeOf((*MockIPKeeper)(nil).GetBlacklist))
}

// GetWhitelist mocks base method
func (m *MockIPKeeper) GetWhitelist() []*net.IPNet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWhitelist")
	ret0, _ := ret[0].([]*net.IPNet)
	return ret0
}

// GetWhitelist indicates an expected call of GetWhitelist
func (mr *MockIPKeeperMockRecorder) GetWhitelist() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWhitelist", reflect.TypeOf((*MockIPKeeper)(nil).GetWhitelist))
}

// Refresh mocks base method
func (m *MockIPKeeper) Refresh(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockIPKeeperMockRecorder) Refresh(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockIPKeeper)(nil).Refresh), ctx)
}

// Watch mocks base method
func (m *MockIPKeeper) Watch(ctx context.Context, interval time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", ctx, interval)
}

// Watch indicates an expected call of Watch
func (mr *MockIPKeeperMockRecorder) Watch(ctx, interval interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockIPKeeper)(nil).Watch), ctx, interval)
}
