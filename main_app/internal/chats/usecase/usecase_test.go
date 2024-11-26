package usecase_test

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"testing"

// 	mock_amqp "git.canopsis.net/canopsis/go-engines/mocks/lib/amqp"
// 	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
// 	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
// 	chatMockRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/repository/mocks"
// 	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/usecase"
// 	messagesMockRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/repository/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/google/uuid"
// )

// func TestAddNewChat(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	ch := mock_amqp.MockChannel{}

// 	chatRepo := chatMockRepo.NewMockChatRepository(ctrl)
// 	messagesRepo := messagesMockRepo.NewMockMessageRepository(ctrl)
// 	usecase := usecase.NewChatUsecase(chatRepo, messagesRepo, ch)

// 	tests := []struct {
// 		name          string
// 		chat          chatModel.ChatDTOInput
// 		prepareMock   func()
// 		prepareCtx    func() context.Context
// 		expectedError error
// 	}{
// 		{
// 			name: "success",
// 			chat: chatModel.ChatDTOInput{
// 				ChatName:   "chat",
// 				ChatType:   "personal",
// 				UsersToAdd: []uuid.UUID{},
// 			},
// 			prepareCtx: func() context.Context {
// 				return context.WithValue(context.Background(), auth.UserKey, auth.User{ID: uuid.New()})
// 			},
// 			prepareMock: func() {
// 				chatRepo.EXPECT().AddUserIntoChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 				chatRepo.EXPECT().CreateNewChat(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.prepareMock()
// 			tt.prepareCtx()

// 			_, err := usecase.AddNewChat(context.Background(), []*http.Cookie{}, tt.chat)

// 			if !errors.Is(err, tt.expectedError) {
// 				t.Errorf("Expected error: '%v', got: '%v'", tt.expectedError, err)
// 			}
// 		})
// 	}
// }
