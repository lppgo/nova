// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/uber/kraken/lib/persistedretry (interfaces: Manager)

// Package mockpersistedretry is a generated GoMock package.
package mockpersistedretry

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	persistedretry "github.com/lppgo/nova/lib/persistedretry"
)

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockManager) Add(arg0 persistedretry.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockManagerMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockManager)(nil).Add), arg0)
}

// Close mocks base method
func (m *MockManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockManager)(nil).Close))
}

// Find mocks base method
func (m *MockManager) Find(arg0 interface{}) ([]persistedretry.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].([]persistedretry.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockManagerMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockManager)(nil).Find), arg0)
}

// SyncExec mocks base method
func (m *MockManager) SyncExec(arg0 persistedretry.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncExec", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncExec indicates an expected call of SyncExec
func (mr *MockManagerMockRecorder) SyncExec(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncExec", reflect.TypeOf((*MockManager)(nil).SyncExec), arg0)
}
