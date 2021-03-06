// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/uber/kraken/origin/blobclient (interfaces: ClusterClient)

// Package mockblobclient is a generated GoMock package.
package mockblobclient

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	core "github.com/lppgo/nova/core"
)

// MockClusterClient is a mock of ClusterClient interface
type MockClusterClient struct {
	ctrl     *gomock.Controller
	recorder *MockClusterClientMockRecorder
}

// MockClusterClientMockRecorder is the mock recorder for MockClusterClient
type MockClusterClientMockRecorder struct {
	mock *MockClusterClient
}

// NewMockClusterClient creates a new mock instance
func NewMockClusterClient(ctrl *gomock.Controller) *MockClusterClient {
	mock := &MockClusterClient{ctrl: ctrl}
	mock.recorder = &MockClusterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClusterClient) EXPECT() *MockClusterClientMockRecorder {
	return m.recorder
}

// DownloadBlob mocks base method
func (m *MockClusterClient) DownloadBlob(arg0 string, arg1 core.Digest, arg2 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadBlob", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadBlob indicates an expected call of DownloadBlob
func (mr *MockClusterClientMockRecorder) DownloadBlob(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadBlob", reflect.TypeOf((*MockClusterClient)(nil).DownloadBlob), arg0, arg1, arg2)
}

// GetMetaInfo mocks base method
func (m *MockClusterClient) GetMetaInfo(arg0 string, arg1 core.Digest) (*core.MetaInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetaInfo", arg0, arg1)
	ret0, _ := ret[0].(*core.MetaInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetaInfo indicates an expected call of GetMetaInfo
func (mr *MockClusterClientMockRecorder) GetMetaInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetaInfo", reflect.TypeOf((*MockClusterClient)(nil).GetMetaInfo), arg0, arg1)
}

// OverwriteMetaInfo mocks base method
func (m *MockClusterClient) OverwriteMetaInfo(arg0 core.Digest, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OverwriteMetaInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OverwriteMetaInfo indicates an expected call of OverwriteMetaInfo
func (mr *MockClusterClientMockRecorder) OverwriteMetaInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OverwriteMetaInfo", reflect.TypeOf((*MockClusterClient)(nil).OverwriteMetaInfo), arg0, arg1)
}

// Owners mocks base method
func (m *MockClusterClient) Owners(arg0 core.Digest) ([]core.PeerContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Owners", arg0)
	ret0, _ := ret[0].([]core.PeerContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Owners indicates an expected call of Owners
func (mr *MockClusterClientMockRecorder) Owners(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Owners", reflect.TypeOf((*MockClusterClient)(nil).Owners), arg0)
}

// ReplicateToRemote mocks base method
func (m *MockClusterClient) ReplicateToRemote(arg0 string, arg1 core.Digest, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplicateToRemote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplicateToRemote indicates an expected call of ReplicateToRemote
func (mr *MockClusterClientMockRecorder) ReplicateToRemote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplicateToRemote", reflect.TypeOf((*MockClusterClient)(nil).ReplicateToRemote), arg0, arg1, arg2)
}

// Stat mocks base method
func (m *MockClusterClient) Stat(arg0 string, arg1 core.Digest) (*core.BlobInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", arg0, arg1)
	ret0, _ := ret[0].(*core.BlobInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat
func (mr *MockClusterClientMockRecorder) Stat(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockClusterClient)(nil).Stat), arg0, arg1)
}

// UploadBlob mocks base method
func (m *MockClusterClient) UploadBlob(arg0 string, arg1 core.Digest, arg2 io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadBlob", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadBlob indicates an expected call of UploadBlob
func (mr *MockClusterClientMockRecorder) UploadBlob(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadBlob", reflect.TypeOf((*MockClusterClient)(nil).UploadBlob), arg0, arg1, arg2)
}
