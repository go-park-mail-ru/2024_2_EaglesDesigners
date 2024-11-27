package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/delivery"
	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	mocks "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddNewChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)
	chatDelivery := delivery.NewChatDelivery(mockService)

	tests := []struct {
		name                 string
		chatData             model.ChatDTOInput
		mockAddNewChatReturn model.ChatDTOOutput
		mockAddNewChatErr    error
		expectedStatusCode   int
	}{
		{
			name: "Successful chat creation",
			chatData: model.ChatDTOInput{
				ChatName:   "Test Chat",
				ChatType:   "group",
				UsersToAdd: []uuid.UUID{uuid.New()},
			},
			mockAddNewChatReturn: model.ChatDTOOutput{
				ChatId:       uuid.New(),
				ChatName:     "Test Chat",
				CountOfUsers: 1,
				ChatType:     "group",
			},
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocking the AddNewChat method
			mockService.EXPECT().
				AddNewChat(gomock.Any(), gomock.Any(), tt.chatData).
				Return(tt.mockAddNewChatReturn, tt.mockAddNewChatErr).
				Times(1)

			// Create a new HTTP request with multipart/form-data
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Write chat_data as a field in the multipart form
			chatDataJSON, _ := json.Marshal(tt.chatData)
			writer.WriteField("chat_data", string(chatDataJSON))

			writer.Close()
			req := httptest.NewRequest(http.MethodPost, "/addchat", body)
			req.Header.Set("Content-Type", writer.FormDataContentType()) // Set correct Content-Type for multipart

			res := httptest.NewRecorder()

			// Call the AddNewChat function
			chatDelivery.AddNewChat(res, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatusCode, res.Code)
		})
	}
}

func TestAddNewChat_InvalidMultipartForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)
	chatDelivery := delivery.NewChatDelivery(mockService)

	req := httptest.NewRequest(http.MethodPost, "/addchat", nil)
	res := httptest.NewRecorder()

	// Call the AddNewChat function without multipart form data
	chatDelivery.AddNewChat(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestGetUserChats(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)

	delivery := delivery.NewChatDelivery(mockService)

	chatList := []model.ChatDTOOutput{

		{ChatId: uuid.New(), ChatName: "Chat 1"},

		{ChatId: uuid.New(), ChatName: "Chat 2"},
	}

	// Настройка мока

	mockService.EXPECT().GetChats(gomock.Any(), gomock.Any()).Return(chatList, nil)

	// Создание тестового HTTP-запроса

	req := httptest.NewRequest(http.MethodGet, "/chats", nil)

	w := httptest.NewRecorder()

	delivery.GetUserChatsHandler(w, req)

	// Проверка ответа

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response model.ChatsDTO

	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, len(chatList), len(response.Chats))

}

func TestDeleteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)
	delivery := delivery.NewChatDelivery(mockService)

	// Создание UUID для пользователя и чата
	user := auth.User{
		ID:       uuid.New(),
		Username: "test",
		Name:     "test",
		Password: "test",
		Version:  1,
	}

	chatID := uuid.New()

	// Настройка мока на успешное удаление чата
	mockService.EXPECT().DeleteChat(gomock.Any(), chatID, user.ID).Return(nil)

	// Создание тестового HTTP-запроса
	req := httptest.NewRequest(http.MethodDelete, "/chat/"+chatID.String()+"/delete", nil)

	// Создание карты для переменных маршрута
	vars := map[string]string{
		"chatId": chatID.String(),
	}

	ctx := context.WithValue(context.Background(), auth.MuxParamsKey, vars)
	ctx = context.WithValue(ctx, auth.UserKey, user)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Вызов метода DeleteChatOrGroup
	delivery.DeleteChatOrGroup(w, req)

	// Проверка ответа
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAddUsersIntoChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)
	chatDelivery := delivery.NewChatDelivery(mockService)

	user := auth.User{
		ID:       uuid.New(),
		Username: "test",
		Name:     "test",
		Password: "test",
		Version:  1,
	}

	chatID := uuid.New().String() // генерируем уникальный идентификатор для чата в строковом формате

	tests := []struct {
		name               string
		usersToAdd         model.AddUsersIntoChatDTO
		mockAddUsersReturn model.AddedUsersIntoChatDTO
		mockAddUsersErr    error
		expectedStatusCode int
	}{
		{
			name: "Successful users addition",
			usersToAdd: model.AddUsersIntoChatDTO{
				UsersId: []uuid.UUID{uuid.New(), uuid.New()},
			},
			mockAddUsersReturn: model.AddedUsersIntoChatDTO{
				AddedUsers:    []uuid.UUID{},
				NotAddedUsers: []uuid.UUID{},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Error decoding JSON",
			usersToAdd: model.AddUsersIntoChatDTO{
				UsersId: []uuid.UUID{},
			},
			mockAddUsersReturn: model.AddedUsersIntoChatDTO{},
			mockAddUsersErr:    assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Invalid chat ID",
			usersToAdd: model.AddUsersIntoChatDTO{
				UsersId: []uuid.UUID{},
			},
			mockAddUsersReturn: model.AddedUsersIntoChatDTO{},
			mockAddUsersErr:    nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocking the AddUsersIntoChatWithCheckPermission method
			mockService.EXPECT().
				AddUsersIntoChatWithCheckPermission(gomock.Any(), tt.usersToAdd.UsersId, chatID).
				Return(tt.mockAddUsersReturn, tt.mockAddUsersErr).
				Times(1)

			// Creating JSON request body
			body, _ := json.Marshal(tt.usersToAdd)
			req := httptest.NewRequest(http.MethodPost, "/chat/"+chatID+"/addusers", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			vars := map[string]string{
				"chatId": chatID,
			}

			ctx := context.WithValue(context.Background(), auth.MuxParamsKey, vars)
			ctx = context.WithValue(ctx, auth.UserKey, user)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			// Call the AddUsersIntoChat function
			chatDelivery.AddUsersIntoChat(res, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatusCode, res.Code)
		})
	}
}

func TestUpdateGroup(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockService := mocks.NewMockChatUsecase(ctrl)

	chatDelivery := delivery.NewChatDelivery(mockService)
	userID := uuid.New() // замените на реальный идентификатор пользователя, если необходимо

	chatID := uuid.New().String()

	user := auth.User{
		ID:       userID,
		Username: "test",
		Name:     "test",
		Password: "test",
		Version:  1,
	}

	tests := []struct {
		name string

		chatUpdate model.ChatUpdate

		mockUpdateReturn model.ChatUpdateOutput

		mockUpdateErr error

		expectedStatusCode int
	}{

		{

			name: "Successful update",

			chatUpdate: model.ChatUpdate{

				ChatName: "Updated Chat Name",

				Avatar: nil, // Здесь может быть nil или указатель на файл

			},

			mockUpdateReturn: model.ChatUpdateOutput{

				ChatName: "Updated Chat Name",

				Avatar: "",
			},

			expectedStatusCode: http.StatusOK,
		},

		{

			name: "Invalid JSON data",

			chatUpdate: model.ChatUpdate{

				ChatName: "Invalid", // Здесь просто валидное имя, но в самом JSON будет неверно

			},

			mockUpdateReturn: model.ChatUpdateOutput{},

			mockUpdateErr: assert.AnError,

			expectedStatusCode: http.StatusInternalServerError,
		},

		{

			name: "No permission error",

			chatUpdate: model.ChatUpdate{

				ChatName: "Another Chat Name",
			},

			mockUpdateReturn: model.ChatUpdateOutput{},

			mockUpdateErr: assert.AnError,

			expectedStatusCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Mocking the UpdateChat method

			mockService.EXPECT().
				UpdateChat(gomock.Any(), chatID, tt.chatUpdate, userID).
				Return(tt.mockUpdateReturn, tt.mockUpdateErr).
				Times(1)

			// Create a multipart/form-data request

			body := &bytes.Buffer{}

			writer := multipart.NewWriter(body)

			// Marshal chatUpdate into JSON and add it to the multipart form

			chatUpdateJSON, _ := json.Marshal(tt.chatUpdate)

			writer.WriteField("chat_data", string(chatUpdateJSON))

			// Optionally add an avatar file, just for simulation

			/*

				if tt.chatUpdate.Avatar != nil {

					part, err := writer.CreateFormFile("avatar", "avatar.png")

					if err != nil {

						t.Fatalf("Failed to create form file: %v", err)

					}

					part.Write([]byte("fake image data")) // Simulated image data

				}

			*/

			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/chat/"+chatID, body)

			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Add user information to the context
			vars := map[string]string{
				"chatId": chatID,
			}

			ctx := context.WithValue(context.Background(), auth.MuxParamsKey, vars)
			ctx = context.WithValue(ctx, auth.UserKey, user)
			req = req.WithContext(ctx)

			req = req.WithContext(ctx)

			res := httptest.NewRecorder()

			// Call the UpdateGroup function

			chatDelivery.UpdateGroup(res, req)

			// Assert the status code

			assert.Equal(t, tt.expectedStatusCode, res.Code)

		})

	}

}
