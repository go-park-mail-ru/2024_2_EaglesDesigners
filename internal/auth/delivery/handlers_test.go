package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/delivery"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	jwtUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_Success(t *testing.T) {
	mockUsecase := &mocks.MockUsecase{
		AuthenticateFunc: func(username, password string) bool {
			return username == "user1" && password == "pass1"
		},
	}
	mockTokenService := &mocks.MockTokenUsecase{
		CreateJWTFunc: func(username string) (string, error) {
			return "mock_token", nil
		},
	}
	handler := delivery.NewDelivery(mockUsecase, mockTokenService)
	reqBody, err := json.Marshal(delivery.AuthCredentials{
		Username: "user1",
		Password: "pass1",
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	result := httptest.NewRecorder()
	handler.LoginHandler(result, request)

	assert.Equal(t, http.StatusCreated, result.Code)
	assert.Contains(t, result.Header().Get("Set-Cookie"), "access_token")

}

func TestRegisterHandler_Success(t *testing.T) {
	mockUsecase := &mocks.MockUsecase{
		RegistrationFunc: func(username, name, password string) error {
			return nil
		},
		GetUserDataFunc: func(username string) (usecase.UserData, error) {
			return usecase.UserData{
				ID:       1,
				Username: username,
				Name:     "name",
			}, nil
		},
	}

	mockTokenService := &mocks.MockTokenUsecase{
		CreateJWTFunc: func(username string) (string, error) {
			return "mock_token", nil
		},
	}
	handler := delivery.NewDelivery(mockUsecase, mockTokenService)
	reqBody, _ := json.Marshal(delivery.RegisterCredentials{
		Username: "killer1994",
		Name:     "Vincent Vega",
		Password: "go_do_a_crime",
	})
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()
	handler.RegisterHandler(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestLogoutHandler_Success(t *testing.T) {
	mockUsecase := &mocks.MockUsecase{}
	mockTokenService := &mocks.MockTokenUsecase{}
	handler := delivery.NewDelivery(mockUsecase, mockTokenService)
	request := httptest.NewRequest(http.MethodPost, "/logout", nil)
	request.AddCookie(&http.Cookie{Name: "access_token", Value: "mock_token"})
	result := httptest.NewRecorder()
	handler.LogoutHandler(result, request)

	assert.Equal(t, http.StatusOK, result.Code)
}

func TestAuthHandler_Success(t *testing.T) {
	mockUsecase := &mocks.MockUsecase{}
	mockTokenService := &mocks.MockTokenUsecase{
		GetUserDataByJWTFunc: func(cookies []*http.Cookie) (jwtUC.UserData, error) {
			return jwtUC.UserData{Username: "user1", Name: "User One"}, nil
		},
	}
	handler := delivery.NewDelivery(mockUsecase, mockTokenService)
	request := httptest.NewRequest(http.MethodGet, "/auth", nil)
	request.AddCookie(&http.Cookie{Name: "access_token", Value: "mock_token"})
	result := httptest.NewRecorder()
	handler.AuthHandler(result, request)

	assert.Equal(t, http.StatusOK, result.Code)

}
