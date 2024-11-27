// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
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

// AddContact mocks base method.
func (m *MockRepository) AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddContact", ctx, contactData)
	ret0, _ := ret[0].(models.ContactDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddContact indicates an expected call of AddContact.
func (mr *MockRepositoryMockRecorder) AddContact(ctx, contactData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddContact", reflect.TypeOf((*MockRepository)(nil).AddContact), ctx, contactData)
}

// DeleteContact mocks base method.
func (m *MockRepository) DeleteContact(ctx context.Context, contactData models.ContactDataDAO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContact", ctx, contactData)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContact indicates an expected call of DeleteContact.
func (mr *MockRepositoryMockRecorder) DeleteContact(ctx, contactData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContact", reflect.TypeOf((*MockRepository)(nil).DeleteContact), ctx, contactData)
}

// GetContacts mocks base method.
func (m *MockRepository) GetContacts(ctx context.Context, username string) ([]models.ContactDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContacts", ctx, username)
	ret0, _ := ret[0].([]models.ContactDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContacts indicates an expected call of GetContacts.
func (mr *MockRepositoryMockRecorder) GetContacts(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContacts", reflect.TypeOf((*MockRepository)(nil).GetContacts), ctx, username)
}

// SearchGlobalUsers mocks base method.
func (m *MockRepository) SearchGlobalUsers(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchGlobalUsers", ctx, userID, keyWord)
	ret0, _ := ret[0].([]models.ContactDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchGlobalUsers indicates an expected call of SearchGlobalUsers.
func (mr *MockRepositoryMockRecorder) SearchGlobalUsers(ctx, userID, keyWord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchGlobalUsers", reflect.TypeOf((*MockRepository)(nil).SearchGlobalUsers), ctx, userID, keyWord)
}

// SearchUserContacts mocks base method.
func (m *MockRepository) SearchUserContacts(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserContacts", ctx, userID, keyWord)
	ret0, _ := ret[0].([]models.ContactDAO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUserContacts indicates an expected call of SearchUserContacts.
func (mr *MockRepositoryMockRecorder) SearchUserContacts(ctx, userID, keyWord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserContacts", reflect.TypeOf((*MockRepository)(nil).SearchUserContacts), ctx, userID, keyWord)
}