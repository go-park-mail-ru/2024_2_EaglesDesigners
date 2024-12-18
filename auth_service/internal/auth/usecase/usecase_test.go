package usecase_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/usecase"
	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/usecase/mocks"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
)

func TestAuthenticate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	// Настройка ожиданий
	mockRepo.EXPECT().GetUserByUsername(ctx, "test").Return(models.UserDAO{
		ID:        uuid.New(),
		Username:  "test",
		Name:      "test",
		Password:  "test",
		Version:   0,
		AvatarURL: new(string),
	}, nil)

	// Вызов метода
	_, err := usecase.Authenticate(ctx, &authv1.AuthRequest{Username: "test", Password: "test"})

	// Проверка результатов
	assert.Nil(t, err)
}

func TestRegistration_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	// Настройка ожиданий
	mockRepo.EXPECT().CreateUser(ctx, "validusername", "testname", gomock.Any()).Return(nil)

	// Вызов метода
	response, err := usecase.Registration(ctx, &authv1.RegistrationRequest{
		Username: "validusername",
		Password: "validpassword123",
		Name:     "testname",
	})

	// Проверка результатов
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Dummy)
}

func TestRegistration_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	// Вызов метода с неправильными данными
	response, err := usecase.Registration(ctx, &authv1.RegistrationRequest{
		Username: "user", // слишком короткий
		Password: "pass", // слишком короткий
		Name:     "",
	})

	// Проверка результатов
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "bad data", err.Error())
}

func TestRegistration_CreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	// Настройка ожиданий для ошибки создания пользователя
	mockRepo.EXPECT().CreateUser(ctx, "validusername", "testname", gomock.Any()).Return(errors.New("database error"))

	// Вызов метода
	response, err := usecase.Registration(ctx, &authv1.RegistrationRequest{
		Username: "validusername",
		Password: "validpassword123",
		Name:     "testname",
	})

	// Проверка результатов
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "database error", err.Error())
}

// func TestIsAuthorized_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := repo.NewMockrepository(ctrl)
// 	mockUsecase := mock_api.NewMockAuth(ctrl)
// 	usecase := usecase.NewUsecase(mockRepo)

// 	ctx := context.Background()
// 	token := "valid.jwt.token"
// 	expectedUser := &authv1.UserJWT{
// 		Username: "testuser",
// 		Version:  1,
// 	}

// 	// Настройка ожиданий
// 	mockUsecase.EXPECT().GetUserByJWT(ctx, gomock.Any()).Return(expectedUser, nil)

// 	// Вызов метода с валидным токеном
// 	response, err := usecase.IsAuthorized(ctx, &authv1.Token{Token: token})

// 	// Проверка результатов
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedUser, response)
// }

func TestIsAuthorized_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	// Вызов метода с невалидным токеном
	response, err := usecase.IsAuthorized(ctx, &authv1.Token{Token: "invalid.token"})

	// Проверка результатов
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid token", err.Error())
}

func TestCreateJWT_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	username := "testuser"
	expectedUser := models.UserDAO{
		ID:       uuid.New(),
		Username: username,
		Name:     "Test User",
		Version:  1,
		Password: "hashed_password", // Предположим, это захешированный пароль
	}

	// Настройка ожиданий
	mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(expectedUser, nil)

	// Вызов метода
	tokenResponse, err := usecase.CreateJWT(ctx, &authv1.CreateJWTRequest{Username: username})

	// Проверка результатов
	assert.Nil(t, err)
	assert.NotNil(t, tokenResponse)
	assert.NotEmpty(t, tokenResponse.Token)
}

func TestCreateJWT_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	usecase := usecase.NewUsecase(mockRepo)
	ctx := context.Background()

	username := "nonexistentuser"

	// Настройка ожиданий
	mockRepo.EXPECT().GetUserByUsername(ctx, username).Return(models.UserDAO{}, errors.New("user not found"))

	// Вызов метода
	tokenResponse, err := usecase.CreateJWT(ctx, &authv1.CreateJWTRequest{Username: username})

	// Проверка результатов
	assert.NotNil(t, err)
	assert.Nil(t, tokenResponse)
	assert.Equal(t, "user not found", err.Error())
}

func TestGetUserByJWT(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := repo.NewMockrepository(ctrl)
	u := usecase.NewUsecase(mockRepo)

	validToken := createJWTToken("testuser")

	tests := []struct {
		name          string
		token         *authv1.Token
		mockSetup     func()
		expectedUser  *authv1.UserJWT
		expectedError error
	}{
		{
			name:  "Valid JWT",
			token: &authv1.Token{Token: validToken},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "testuser").
					Return(models.UserDAO{ID: uuid.MustParse("27d4a2da-60a0-4a89-bd68-282526ee108e"), Username: "testuser", Name: "Test User", Password: "password"}, nil)
			},
			expectedUser: &authv1.UserJWT{
				ID:       "27d4a2da-60a0-4a89-bd68-282526ee108e",
				Username: "testuser",
				Name:     "Test User",
				Password: "password",
			},

			expectedError: nil,
		},

		{
			name:  "Invalid JWT",
			token: &authv1.Token{Token: "invalid.jwt.token"},
			mockSetup: func() {
			},
			expectedUser:  nil,
			expectedError: errors.New("невалидный jwt token"),
		},
		{
			name:  "User not found",
			token: &authv1.Token{Token: validToken},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "testuser").
					Return(models.UserDAO{}, errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := u.GetUserByJWT(context.Background(), tt.token)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}

func createJWTToken(username string) string {
	payload := models.Payload{
		Sub: username,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	token := fmt.Sprintf("header.%s.signature", base64.RawURLEncoding.EncodeToString(payloadBytes))

	return token
}
