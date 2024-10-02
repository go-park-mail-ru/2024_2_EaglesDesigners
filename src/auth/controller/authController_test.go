package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/controller"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_Success(t *testing.T) {
	mockAuthService := &mocks.MockAuthService{
		AuthenticateFunc: func(username, password string) bool {
			return username == "user1" && password == "pass1"
		},
	}
	mockTokenService := &mocks.MockTokenService{
		CreateJWTFunc: func(username string) (string, error) {
			return "mock_token", nil
		},
	}
	controller := controller.NewAuthController(mockAuthService, mockTokenService)
	reqBody, _ := json.Marshal(utils.AuthCredentials{
		Username: "user1",
		Password: "pass1",
	})
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	result := httptest.NewRecorder()
	controller.LoginHandler(result, request)

	assert.Equal(t, http.StatusCreated, result.Code)
	assert.Contains(t, result.Header().Get("Set-Cookie"), "access_token")

}

func TestRegisterHandler_Success(t *testing.T) {
	mockAuthService := &mocks.MockAuthService{
		RegistrationFunc: func(username, name, password string) error {
			return nil
		},
		GetUserDataFunc: func(username string) (utils.UserData, error) {
			return utils.UserData{
				ID:       1,
				Username: username,
				Name:     "name",
			}, nil
		},
	}

	mockTokenService := &mocks.MockTokenService{
		CreateJWTFunc: func(username string) (string, error) {
			return "mock_token", nil
		},
	}
	controller := controller.NewAuthController(mockAuthService, mockTokenService)
	reqBody, _ := json.Marshal(utils.RegisterCredentials{
		Username: "killer1994",
		Name:     "Vincent Vega",
		Password: "go_do_a_crime",
	})
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()
	controller.RegisterHandler(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestLogoutHandler_Success(t *testing.T) {
	mockAuthService := &mocks.MockAuthService{}
	mockTokenService := &mocks.MockTokenService{}
	controller := controller.NewAuthController(mockAuthService, mockTokenService)
	request := httptest.NewRequest(http.MethodPost, "/logout", nil)
	request.AddCookie(&http.Cookie{Name: "access_token", Value: "mock_token"})
	result := httptest.NewRecorder()
	controller.LogoutHandler(result, request)

	assert.Equal(t, http.StatusOK, result.Code)
}

func TestAuthHandler_Success(t *testing.T) {
	mockAuthService := &mocks.MockAuthService{}
	mockTokenService := &mocks.MockTokenService{
		GetUserDataByJWTFunc: func(cookies []*http.Cookie) (utils.UserData, error) {
			return utils.UserData{Username: "user1", Name: "User One"}, nil
		},
	}
	controller := controller.NewAuthController(mockAuthService, mockTokenService)
	request := httptest.NewRequest(http.MethodGet, "/auth", nil)
	request.AddCookie(&http.Cookie{Name: "access_token", Value: "mock_token"})
	result := httptest.NewRecorder()
	controller.AuthHandler(result, request)

	assert.Equal(t, http.StatusOK, result.Code)

}
