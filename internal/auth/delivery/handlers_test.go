package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/delivery"
	mock_delivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/delivery/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_delivery.NewMockusecase(ctrl)
	mockToken := mock_delivery.NewMocktoken(ctrl)
	delivery := delivery.NewDelivery(mockUsecase, mockToken)

	handler := http.HandlerFunc(delivery.LoginHandler)

	tests := []struct {
		name            string
		username        string
		password        string
		mockAuthReturn  bool
		mockTokenReturn string
		expectedStatus  int
	}{
		{
			name:            "Successful Authentication",
			username:        "user11",
			password:        "validPassword",
			mockAuthReturn:  true,
			mockTokenReturn: "someJWTToken",
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Failed Authentication - Wrong Password",
			username:        "user22",
			password:        "wrongPassword",
			mockAuthReturn:  false,
			mockTokenReturn: "",
			expectedStatus:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.EXPECT().Authenticate(gomock.Any(), tt.username, tt.password).Return(tt.mockAuthReturn)

			if tt.mockAuthReturn {
				mockToken.EXPECT().CreateJWT(gomock.Any(), tt.username).Return(tt.mockTokenReturn, nil)
			}

			reqBody, _ := json.Marshal(models.AuthReqDTO{
				Username: tt.username,
				Password: tt.password,
			})

			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock_delivery.NewMockusecase(ctrl)
	mockToken := mock_delivery.NewMocktoken(ctrl)
	delivery := delivery.NewDelivery(mockUsecase, mockToken)

	handler := http.HandlerFunc(delivery.RegisterHandler)

	tests := []struct {
		name           string
		username       string
		nickname       string
		password       string
		mockReturn     error
		expectedStatus int
	}{
		{
			name:           "Successful Registration",
			username:       "newUser",
			nickname:       "New User",
			password:       "newPassword123",
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failed Registration - User Exists",
			username:       "existingUser",
			nickname:       "Existing User",
			password:       "existingPassword123",
			mockReturn:     models.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.EXPECT().Registration(gomock.Any(), tt.username, tt.nickname, tt.password).Return(tt.mockReturn)

			reqBody, _ := json.Marshal(models.RegisterReqDTO{
				Username: tt.username,
				Name:     tt.nickname,
				Password: tt.password,
			})
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
