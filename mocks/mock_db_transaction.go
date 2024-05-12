// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/itohin/gophkeeper/internal/server/usecases/auth (interfaces: DBTransactionManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDBTransactionManager is a mock of DBTransactionManager interface.
type MockDBTransactionManager struct {
	ctrl     *gomock.Controller
	recorder *MockDBTransactionManagerMockRecorder
}

// MockDBTransactionManagerMockRecorder is the mock recorder for MockDBTransactionManager.
type MockDBTransactionManagerMockRecorder struct {
	mock *MockDBTransactionManager
}

// NewMockDBTransactionManager creates a new mock instance.
func NewMockDBTransactionManager(ctrl *gomock.Controller) *MockDBTransactionManager {
	mock := &MockDBTransactionManager{ctrl: ctrl}
	mock.recorder = &MockDBTransactionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBTransactionManager) EXPECT() *MockDBTransactionManagerMockRecorder {
	return m.recorder
}

// Transaction mocks base method.
func (m *MockDBTransactionManager) Transaction(arg0 context.Context, arg1 func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockDBTransactionManagerMockRecorder) Transaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockDBTransactionManager)(nil).Transaction), arg0, arg1)
}
