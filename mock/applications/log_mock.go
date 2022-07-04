// Code generated by MockGen. DO NOT EDIT.
// Source: applications/log.go

// Package applications_mock is a generated GoMock package.
package applications_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogService is a mock of LogService interface.
type MockLogService struct {
	ctrl     *gomock.Controller
	recorder *MockLogServiceMockRecorder
}

// MockLogServiceMockRecorder is the mock recorder for MockLogService.
type MockLogServiceMockRecorder struct {
	mock *MockLogService
}

// NewMockLogService creates a new mock instance.
func NewMockLogService(ctrl *gomock.Controller) *MockLogService {
	mock := &MockLogService{ctrl: ctrl}
	mock.recorder = &MockLogServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogService) EXPECT() *MockLogServiceMockRecorder {
	return m.recorder
}

// Critical mocks base method.
func (m *MockLogService) Critical(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Critical", varargs...)
}

// Critical indicates an expected call of Critical.
func (mr *MockLogServiceMockRecorder) Critical(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Critical", reflect.TypeOf((*MockLogService)(nil).Critical), args...)
}

// Debug mocks base method.
func (m *MockLogService) Debug(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLogServiceMockRecorder) Debug(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogService)(nil).Debug), args...)
}

// Info mocks base method.
func (m *MockLogService) Info(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLogServiceMockRecorder) Info(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogService)(nil).Info), args...)
}
