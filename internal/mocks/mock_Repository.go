// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pavlegich/scripts-hub/internal/repository (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/pavlegich/scripts-hub/internal/entities"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AppendCommandOutputByName mocks base method.
func (m *MockRepository) AppendCommandOutputByName(arg0 context.Context, arg1 *entities.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendCommandOutputByName", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendCommandOutputByName indicates an expected call of AppendCommandOutputByName.
func (mr *MockRepositoryMockRecorder) AppendCommandOutputByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendCommandOutputByName", reflect.TypeOf((*MockRepository)(nil).AppendCommandOutputByName), arg0, arg1)
}

// CreateCommand mocks base method.
func (m *MockRepository) CreateCommand(arg0 context.Context, arg1 *entities.Command) (*entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommand", arg0, arg1)
	ret0, _ := ret[0].(*entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommand indicates an expected call of CreateCommand.
func (mr *MockRepositoryMockRecorder) CreateCommand(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommand", reflect.TypeOf((*MockRepository)(nil).CreateCommand), arg0, arg1)
}

// DeleteCommandByName mocks base method.
func (m *MockRepository) DeleteCommandByName(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommandByName", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommandByName indicates an expected call of DeleteCommandByName.
func (mr *MockRepositoryMockRecorder) DeleteCommandByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommandByName", reflect.TypeOf((*MockRepository)(nil).DeleteCommandByName), arg0, arg1)
}

// GetAllCommands mocks base method.
func (m *MockRepository) GetAllCommands(arg0 context.Context) ([]*entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCommands", arg0)
	ret0, _ := ret[0].([]*entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCommands indicates an expected call of GetAllCommands.
func (mr *MockRepositoryMockRecorder) GetAllCommands(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCommands", reflect.TypeOf((*MockRepository)(nil).GetAllCommands), arg0)
}

// GetCommandByName mocks base method.
func (m *MockRepository) GetCommandByName(arg0 context.Context, arg1 string) (*entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommandByName", arg0, arg1)
	ret0, _ := ret[0].(*entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommandByName indicates an expected call of GetCommandByName.
func (mr *MockRepositoryMockRecorder) GetCommandByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommandByName", reflect.TypeOf((*MockRepository)(nil).GetCommandByName), arg0, arg1)
}
