// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pavlegich/scripts-hub/internal/service/command (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/pavlegich/scripts-hub/internal/entities"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AppendOutput mocks base method.
func (m *MockService) AppendOutput(arg0 context.Context, arg1 *entities.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendOutput", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendOutput indicates an expected call of AppendOutput.
func (mr *MockServiceMockRecorder) AppendOutput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendOutput", reflect.TypeOf((*MockService)(nil).AppendOutput), arg0, arg1)
}

// Create mocks base method.
func (m *MockService) Create(arg0 context.Context, arg1 *entities.Command) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockService) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), arg0, arg1)
}

// List mocks base method.
func (m *MockService) List(arg0 context.Context) ([]*entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), arg0)
}

// Unload mocks base method.
func (m *MockService) Unload(arg0 context.Context, arg1 string) (*entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unload", arg0, arg1)
	ret0, _ := ret[0].(*entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unload indicates an expected call of Unload.
func (mr *MockServiceMockRecorder) Unload(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unload", reflect.TypeOf((*MockService)(nil).Unload), arg0, arg1)
}
