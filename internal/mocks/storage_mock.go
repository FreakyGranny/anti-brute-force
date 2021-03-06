// Code generated by MockGen. DO NOT EDIT.
// Source: entity.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	storage "github.com/FreakyGranny/anti-brute-force/internal/storage"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockWriteStorage is a mock of WriteStorage interface
type MockWriteStorage struct {
	ctrl     *gomock.Controller
	recorder *MockWriteStorageMockRecorder
}

// MockWriteStorageMockRecorder is the mock recorder for MockWriteStorage
type MockWriteStorageMockRecorder struct {
	mock *MockWriteStorage
}

// NewMockWriteStorage creates a new mock instance
func NewMockWriteStorage(ctrl *gomock.Controller) *MockWriteStorage {
	mock := &MockWriteStorage{ctrl: ctrl}
	mock.recorder = &MockWriteStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWriteStorage) EXPECT() *MockWriteStorageMockRecorder {
	return m.recorder
}

// AddToWhiteList mocks base method
func (m *MockWriteStorage) AddToWhiteList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWhiteList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWhiteList indicates an expected call of AddToWhiteList
func (mr *MockWriteStorageMockRecorder) AddToWhiteList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWhiteList", reflect.TypeOf((*MockWriteStorage)(nil).AddToWhiteList), ctx, ip, mask)
}

// AddToBlackList mocks base method
func (m *MockWriteStorage) AddToBlackList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlackList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlackList indicates an expected call of AddToBlackList
func (mr *MockWriteStorageMockRecorder) AddToBlackList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlackList", reflect.TypeOf((*MockWriteStorage)(nil).AddToBlackList), ctx, ip, mask)
}

// RemoveFromWhiteList mocks base method
func (m *MockWriteStorage) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromWhiteList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromWhiteList indicates an expected call of RemoveFromWhiteList
func (mr *MockWriteStorageMockRecorder) RemoveFromWhiteList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromWhiteList", reflect.TypeOf((*MockWriteStorage)(nil).RemoveFromWhiteList), ctx, ip, mask)
}

// RemoveFromBlackList mocks base method
func (m *MockWriteStorage) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromBlackList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromBlackList indicates an expected call of RemoveFromBlackList
func (mr *MockWriteStorageMockRecorder) RemoveFromBlackList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBlackList", reflect.TypeOf((*MockWriteStorage)(nil).RemoveFromBlackList), ctx, ip, mask)
}

// MockReadStorage is a mock of ReadStorage interface
type MockReadStorage struct {
	ctrl     *gomock.Controller
	recorder *MockReadStorageMockRecorder
}

// MockReadStorageMockRecorder is the mock recorder for MockReadStorage
type MockReadStorageMockRecorder struct {
	mock *MockReadStorage
}

// NewMockReadStorage creates a new mock instance
func NewMockReadStorage(ctrl *gomock.Controller) *MockReadStorage {
	mock := &MockReadStorage{ctrl: ctrl}
	mock.recorder = &MockReadStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReadStorage) EXPECT() *MockReadStorageMockRecorder {
	return m.recorder
}

// GetBlackList mocks base method
func (m *MockReadStorage) GetBlackList(ctx context.Context) ([]*storage.IPNet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlackList", ctx)
	ret0, _ := ret[0].([]*storage.IPNet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlackList indicates an expected call of GetBlackList
func (mr *MockReadStorageMockRecorder) GetBlackList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlackList", reflect.TypeOf((*MockReadStorage)(nil).GetBlackList), ctx)
}

// GetWhiteList mocks base method
func (m *MockReadStorage) GetWhiteList(ctx context.Context) ([]*storage.IPNet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWhiteList", ctx)
	ret0, _ := ret[0].([]*storage.IPNet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWhiteList indicates an expected call of GetWhiteList
func (mr *MockReadStorageMockRecorder) GetWhiteList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWhiteList", reflect.TypeOf((*MockReadStorage)(nil).GetWhiteList), ctx)
}
