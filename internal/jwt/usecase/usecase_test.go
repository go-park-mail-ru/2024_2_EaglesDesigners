package usecase_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	mocks "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	u := usecase.NewUsecase(mockRepo)

	tests := []struct {
		name        string
		cookies     []*http.Cookie
		mockGetUser func()
		expectedErr error
	}{
		{
			name: "valid token",
			cookies: []*http.Cookie{
				{Name: "access_token", Value: generateValidJWT("testuser")},
			},
			mockGetUser: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "testuser").Return(models.UserDAO{Username: "testuser", Version: 1}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "invalid token",
			cookies: []*http.Cookie{
				{Name: "access_token", Value: "invalid.jwt.token"},
			},
			mockGetUser: func() {},
			expectedErr: errors.New("токен невалиден"), // Убедитесь, что сообщение совпадает.
		},
		// Добавьте дополнительные тесты для других сценариев.
		// ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockGetUser != nil {
				tt.mockGetUser()
			}
			_, err := u.IsAuthorized(context.Background(), tt.cookies)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func generateValidJWT(username string) string {
	payload := usecase.Payload{
		Sub:     username,
		Name:    "Test User",
		ID:      uuid.New(),
		Version: 1,
		Exp:     time.Now().Add(time.Hour * 24).Unix(),
	}

	header := usecase.Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	h, _ := json.Marshal(header)
	p, _ := json.Marshal(payload)

	headerEncoded := base64.RawURLEncoding.EncodeToString(h)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(p)

	sig, _ := usecase.GeneratorJWT(headerEncoded, payloadEncoded, usecase.GenerateJWTSecret())
	return sig
}

func TestGetUserDataByJWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	u := usecase.NewUsecase(mockRepo)

	validJWT := generateValidJWT("testuser")

	tests := []struct {
		name             string
		cookies          []*http.Cookie
		mockGetUser      func()
		expectedUserData usecase.UserData
		expectedErr      error
	}{
		{
			name: "valid token",
			cookies: []*http.Cookie{
				{Name: "access_token", Value: validJWT},
			},
			mockGetUser: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "testuser").Return(models.UserDAO{
					Username: "testuser",
					Name:     "Test User",
					ID:       uuid.New(),
					Version:  1,
				}, nil)
			},
			expectedUserData: usecase.UserData{
				Username: "testuser",
				Name:     "Test User",
			},
			expectedErr: nil,
		},
		{
			name: "invalid token",
			cookies: []*http.Cookie{
				{Name: "access_token", Value: "invalid.jwt.token"},
			},
			mockGetUser: func() {},
			expectedErr: errors.New("невалидный jwt token"),
		},
		{
			name: "user not found",
			cookies: []*http.Cookie{
				{Name: "access_token", Value: validJWT},
			},
			mockGetUser: func() {
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "testuser").Return(models.UserDAO{}, errors.New("user not found"))
			},
			expectedErr: errors.New("user not found"),
		},
		// Additional cases can be added if necessary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockGetUser != nil {
				tt.mockGetUser()
			}
			userData, err := u.GetUserDataByJWT(tt.cookies)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserData, userData)
			}
		})
	}
}
