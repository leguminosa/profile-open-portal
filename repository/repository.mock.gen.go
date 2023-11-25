// Code generated by MockGen. DO NOT EDIT.
// Source: repository/repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/leguminosa/profile-open-portal/entity"
)

// MockUserRepositoryInterface is a mock of UserRepositoryInterface interface.
type MockUserRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryInterfaceMockRecorder
}

// MockUserRepositoryInterfaceMockRecorder is the mock recorder for MockUserRepositoryInterface.
type MockUserRepositoryInterfaceMockRecorder struct {
	mock *MockUserRepositoryInterface
}

// NewMockUserRepositoryInterface creates a new mock instance.
func NewMockUserRepositoryInterface(ctrl *gomock.Controller) *MockUserRepositoryInterface {
	mock := &MockUserRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepositoryInterface) EXPECT() *MockUserRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetUserByPhoneNumber mocks base method.
func (m *MockUserRepositoryInterface) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByPhoneNumber", ctx, phoneNumber)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByPhoneNumber indicates an expected call of GetUserByPhoneNumber.
func (mr *MockUserRepositoryInterfaceMockRecorder) GetUserByPhoneNumber(ctx, phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByPhoneNumber", reflect.TypeOf((*MockUserRepositoryInterface)(nil).GetUserByPhoneNumber), ctx, phoneNumber)
}

// InsertUser mocks base method.
func (m *MockUserRepositoryInterface) InsertUser(ctx context.Context, user *entity.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockUserRepositoryInterfaceMockRecorder) InsertUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockUserRepositoryInterface)(nil).InsertUser), ctx, user)
}
