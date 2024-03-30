// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/itohin/gophkeeper/internal/server/usecases/auth (interfaces: SessionsStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/itohin/gophkeeper/internal/server/entities"
)

// MockSessionsStorage is a mock of SessionsStorage interface.
type MockSessionsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsStorageMockRecorder
}

// MockSessionsStorageMockRecorder is the mock recorder for MockSessionsStorage.
type MockSessionsStorageMockRecorder struct {
	mock *MockSessionsStorage
}

// NewMockSessionsStorage creates a new mock instance.
func NewMockSessionsStorage(ctrl *gomock.Controller) *MockSessionsStorage {
	mock := &MockSessionsStorage{ctrl: ctrl}
	mock.recorder = &MockSessionsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionsStorage) EXPECT() *MockSessionsStorageMockRecorder {
	return m.recorder
}

// DeleteByID mocks base method.
func (m *MockSessionsStorage) DeleteByID(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockSessionsStorageMockRecorder) DeleteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockSessionsStorage)(nil).DeleteByID), arg0, arg1)
}

// DeleteByUserAndFingerPrint mocks base method.
func (m *MockSessionsStorage) DeleteByUserAndFingerPrint(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserAndFingerPrint", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserAndFingerPrint indicates an expected call of DeleteByUserAndFingerPrint.
func (mr *MockSessionsStorageMockRecorder) DeleteByUserAndFingerPrint(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserAndFingerPrint", reflect.TypeOf((*MockSessionsStorage)(nil).DeleteByUserAndFingerPrint), arg0, arg1, arg2)
}

// FindByFingerPrint mocks base method.
func (m *MockSessionsStorage) FindByFingerPrint(arg0 context.Context, arg1, arg2 string) (*entities.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByFingerPrint", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entities.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByFingerPrint indicates an expected call of FindByFingerPrint.
func (mr *MockSessionsStorageMockRecorder) FindByFingerPrint(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByFingerPrint", reflect.TypeOf((*MockSessionsStorage)(nil).FindByFingerPrint), arg0, arg1, arg2)
}

// FindByID mocks base method.
func (m *MockSessionsStorage) FindByID(arg0 context.Context, arg1 string) (*entities.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*entities.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockSessionsStorageMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockSessionsStorage)(nil).FindByID), arg0, arg1)
}

// Save mocks base method.
func (m *MockSessionsStorage) Save(arg0 context.Context, arg1 entities.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSessionsStorageMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessionsStorage)(nil).Save), arg0, arg1)
}
