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
			mockTokenReturn: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTEiLCJpZCI6MTIzLCJleHBhbGl0eSI6MTY0MjcwNDYyMn0.kZex8C1HNV8x_XHg5gGKGh7x8ZgghIFuBFlmQU6-F-o",
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "Failed Authentication - Wrong Password",
			username:        "user22",
			password:        "wrongPassword",
			mockAuthReturn:  false,
			mockTokenReturn: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTEiLCJpZCI6MTIzLCJleHBhbGl0eSI6MTY0MjcwNDYyMn0.kZex8C1HNV8x_XHg5gGKGh7x8ZgghIFuBFlmQU6-F-o",
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
	tests := []struct {
		name           string
		username       string
		nickname       string
		password       string
		mockReturn     error
		jwtReturn      string
		expectedStatus int
	}{
		{
			name:           "Successful Registration",
			username:       "user11",
			nickname:       "New User",
			password:       "12345678",
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_delivery.NewMockusecase(ctrl)
			mockToken := mock_delivery.NewMocktoken(ctrl)
			delivery := delivery.NewDelivery(mockUsecase, mockToken)

			handler := http.HandlerFunc(delivery.RegisterHandler)

			mockUsecase.EXPECT().Registration(gomock.Any(), tt.username, tt.nickname, tt.password).Return(tt.mockReturn)
			mockToken.EXPECT().CreateJWT(gomock.Any(), tt.username).Return()

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
