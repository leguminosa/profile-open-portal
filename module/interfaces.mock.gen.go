// Code generated by MockGen. DO NOT EDIT.
// Source: module/interfaces.go

// Package module is a generated GoMock package.
package module

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/leguminosa/profile-open-portal/entity"
)

// MockUserModuleInterface is a mock of UserModuleInterface interface.
type MockUserModuleInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserModuleInterfaceMockRecorder
}

// MockUserModuleInterfaceMockRecorder is the mock recorder for MockUserModuleInterface.
type MockUserModuleInterfaceMockRecorder struct {
	mock *MockUserModuleInterface
}

// NewMockUserModuleInterface creates a new mock instance.
func NewMockUserModuleInterface(ctrl *gomock.Controller) *MockUserModuleInterface {
	mock := &MockUserModuleInterface{ctrl: ctrl}
	mock.recorder = &MockUserModuleInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserModuleInterface) EXPECT() *MockUserModuleInterfaceMockRecorder {
	return m.recorder
}

// GetProfile mocks base method.
func (m *MockUserModuleInterface) GetProfile(ctx context.Context, userID int) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", ctx, userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockUserModuleInterfaceMockRecorder) GetProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockUserModuleInterface)(nil).GetProfile), ctx, userID)
}

// Login mocks base method.
func (m *MockUserModuleInterface) Login(ctx context.Context, user *entity.User) (entity.LoginModuleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, user)
	ret0, _ := ret[0].(entity.LoginModuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserModuleInterfaceMockRecorder) Login(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserModuleInterface)(nil).Login), ctx, user)
}

// Register mocks base method.
func (m *MockUserModuleInterface) Register(ctx context.Context, user *entity.User) (entity.RegisterModuleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, user)
	ret0, _ := ret[0].(entity.RegisterModuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserModuleInterfaceMockRecorder) Register(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserModuleInterface)(nil).Register), ctx, user)
}

// UpdateProfile mocks base method.
func (m *MockUserModuleInterface) UpdateProfile(ctx context.Context, user *entity.User) (entity.UpdateProfileModuleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", ctx, user)
	ret0, _ := ret[0].(entity.UpdateProfileModuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserModuleInterfaceMockRecorder) UpdateProfile(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserModuleInterface)(nil).UpdateProfile), ctx, user)
}