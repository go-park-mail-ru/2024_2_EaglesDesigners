// Code generated by MockGen. DO NOT EDIT.
// Source: usecase_interface.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	http "net/http"
	reflect "reflect"

	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockChatUsecase is a mock of ChatUsecase interface.
type MockChatUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockChatUsecaseMockRecorder
}

// MockChatUsecaseMockRecorder is the mock recorder for MockChatUsecase.
type MockChatUsecaseMockRecorder struct {
	mock *MockChatUsecase
}

// NewMockChatUsecase creates a new mock instance.
func NewMockChatUsecase(ctrl *gomock.Controller) *MockChatUsecase {
	mock := &MockChatUsecase{ctrl: ctrl}
	mock.recorder = &MockChatUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatUsecase) EXPECT() *MockChatUsecaseMockRecorder {
	return m.recorder
}

// AddBranch mocks base method.
func (m *MockChatUsecase) AddBranch(ctx context.Context, chatId, messageID, userId uuid.UUID) (model.AddBranch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBranch", ctx, chatId, messageID, userId)
	ret0, _ := ret[0].(model.AddBranch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBranch indicates an expected call of AddBranch.
func (mr *MockChatUsecaseMockRecorder) AddBranch(ctx, chatId, messageID, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBranch", reflect.TypeOf((*MockChatUsecase)(nil).AddBranch), ctx, chatId, messageID, userId)
}

// AddNewChat mocks base method.
func (m *MockChatUsecase) AddNewChat(ctx context.Context, cookie []*http.Cookie, chat model.ChatDTOInput) (model.ChatDTOOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewChat", ctx, cookie, chat)
	ret0, _ := ret[0].(model.ChatDTOOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNewChat indicates an expected call of AddNewChat.
func (mr *MockChatUsecaseMockRecorder) AddNewChat(ctx, cookie, chat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewChat", reflect.TypeOf((*MockChatUsecase)(nil).AddNewChat), ctx, cookie, chat)
}

// AddUsersIntoChatWithCheckPermission mocks base method.
func (m *MockChatUsecase) AddUsersIntoChatWithCheckPermission(ctx context.Context, user_ids []uuid.UUID, chat_id uuid.UUID) (model.AddedUsersIntoChatDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUsersIntoChatWithCheckPermission", ctx, user_ids, chat_id)
	ret0, _ := ret[0].(model.AddedUsersIntoChatDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUsersIntoChatWithCheckPermission indicates an expected call of AddUsersIntoChatWithCheckPermission.
func (mr *MockChatUsecaseMockRecorder) AddUsersIntoChatWithCheckPermission(ctx, user_ids, chat_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUsersIntoChatWithCheckPermission", reflect.TypeOf((*MockChatUsecase)(nil).AddUsersIntoChatWithCheckPermission), ctx, user_ids, chat_id)
}

// DeleteChat mocks base method.
func (m *MockChatUsecase) DeleteChat(ctx context.Context, chatId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChat", ctx, chatId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChat indicates an expected call of DeleteChat.
func (mr *MockChatUsecaseMockRecorder) DeleteChat(ctx, chatId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChat", reflect.TypeOf((*MockChatUsecase)(nil).DeleteChat), ctx, chatId, userId)
}

// DeleteUsersFromChat mocks base method.
func (m *MockChatUsecase) DeleteUsersFromChat(ctx context.Context, userID, chatId uuid.UUID, usertToDelete model.DeleteUsersFromChatDTO) (model.DeletdeUsersFromChatDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUsersFromChat", ctx, userID, chatId, usertToDelete)
	ret0, _ := ret[0].(model.DeletdeUsersFromChatDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUsersFromChat indicates an expected call of DeleteUsersFromChat.
func (mr *MockChatUsecaseMockRecorder) DeleteUsersFromChat(ctx, userID, chatId, usertToDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUsersFromChat", reflect.TypeOf((*MockChatUsecase)(nil).DeleteUsersFromChat), ctx, userID, chatId, usertToDelete)
}

// GetChatInfo mocks base method.
func (m *MockChatUsecase) GetChatInfo(ctx context.Context, chatId, userId uuid.UUID) (model.ChatInfoDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatInfo", ctx, chatId, userId)
	ret0, _ := ret[0].(model.ChatInfoDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatInfo indicates an expected call of GetChatInfo.
func (mr *MockChatUsecaseMockRecorder) GetChatInfo(ctx, chatId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatInfo", reflect.TypeOf((*MockChatUsecase)(nil).GetChatInfo), ctx, chatId, userId)
}

// GetChats mocks base method.
func (m *MockChatUsecase) GetChats(ctx context.Context, cookie []*http.Cookie) ([]model.ChatDTOOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChats", ctx, cookie)
	ret0, _ := ret[0].([]model.ChatDTOOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChats indicates an expected call of GetChats.
func (mr *MockChatUsecaseMockRecorder) GetChats(ctx, cookie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChats", reflect.TypeOf((*MockChatUsecase)(nil).GetChats), ctx, cookie)
}

// GetUserChats mocks base method.
func (m *MockChatUsecase) GetUserChats(ctx context.Context, userId string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserChats", ctx, userId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserChats indicates an expected call of GetUserChats.
func (mr *MockChatUsecaseMockRecorder) GetUserChats(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserChats", reflect.TypeOf((*MockChatUsecase)(nil).GetUserChats), ctx, userId)
}

// GetUsersFromChat mocks base method.
func (m *MockChatUsecase) GetUsersFromChat(ctx context.Context, chatId string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersFromChat", ctx, chatId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersFromChat indicates an expected call of GetUsersFromChat.
func (mr *MockChatUsecaseMockRecorder) GetUsersFromChat(ctx, chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersFromChat", reflect.TypeOf((*MockChatUsecase)(nil).GetUsersFromChat), ctx, chatId)
}

// SearchChats mocks base method.
func (m *MockChatUsecase) SearchChats(ctx context.Context, userID uuid.UUID, keyWord string) (model.SearchChatsDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchChats", ctx, userID, keyWord)
	ret0, _ := ret[0].(model.SearchChatsDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchChats indicates an expected call of SearchChats.
func (mr *MockChatUsecaseMockRecorder) SearchChats(ctx, userID, keyWord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchChats", reflect.TypeOf((*MockChatUsecase)(nil).SearchChats), ctx, userID, keyWord)
}

// UpdateChat mocks base method.
func (m *MockChatUsecase) UpdateChat(ctx context.Context, chatId uuid.UUID, chatUpdate model.ChatUpdate, userId uuid.UUID) (model.ChatUpdateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChat", ctx, chatId, chatUpdate, userId)
	ret0, _ := ret[0].(model.ChatUpdateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChat indicates an expected call of UpdateChat.
func (mr *MockChatUsecaseMockRecorder) UpdateChat(ctx, chatId, chatUpdate, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChat", reflect.TypeOf((*MockChatUsecase)(nil).UpdateChat), ctx, chatId, chatUpdate, userId)
}

// UserLeaveChat mocks base method.
func (m *MockChatUsecase) UserLeaveChat(ctx context.Context, userId, chatId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserLeaveChat", ctx, userId, chatId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserLeaveChat indicates an expected call of UserLeaveChat.
func (mr *MockChatUsecaseMockRecorder) UserLeaveChat(ctx, userId, chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLeaveChat", reflect.TypeOf((*MockChatUsecase)(nil).UserLeaveChat), ctx, userId, chatId)
}