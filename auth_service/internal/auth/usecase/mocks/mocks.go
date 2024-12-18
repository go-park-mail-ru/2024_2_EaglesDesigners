// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/models"
	gomock "github.com/golang/mock/gomock"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *Mockrepository) CreateUser(ctx context.Context, username, name, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, username, name, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockrepositoryMockRecorder) CreateUser(ctx, username, name, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*Mockrepository)(nil).CreateUser), ctx, username, name, password)
}

// GetUserByUsername mocks base method.
func (m *Mockrepository) GetUserByUsername(ctx context.Context, username string) (models.UserDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(models.UserDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockrepositoryMockRecorder) GetUserByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*Mockrepository)(nil).GetUserByUsername), ctx, username)
}
