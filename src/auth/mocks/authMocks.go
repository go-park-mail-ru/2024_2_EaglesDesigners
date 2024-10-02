package mocks

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/utils"
)

type MockAuthService struct {
	AuthenticateFunc func(username, password string) bool
	RegistrationFunc func(username, name, password string) error
	GetUserDataFunc  func(username string) (utils.UserData, error)
}

func (m *MockAuthService) Authenticate(username, password string) bool {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(username, password)
	}
	return false
}

func (m *MockAuthService) Registration(username, name, password string) error {
	if m.RegistrationFunc != nil {
		return m.RegistrationFunc(username, name, password)
	}
	return nil
}

func (m *MockAuthService) GetUserDataByUsername(username string) (utils.UserData, error) {
	if m.GetUserDataFunc != nil {
		return m.GetUserDataFunc(username)
	}
	return utils.UserData{}, nil
}

type MockTokenService struct {
	CreateJWTFunc        func(username string) (string, error)
	GetUserDataByJWTFunc func(cookies []*http.Cookie) (utils.UserData, error)
	GetUserByJWTFunc     func(cookies []*http.Cookie) (model.User, error)
	IsAuthorizedFunc     func(cookies []*http.Cookie) error
}

func (m *MockTokenService) CreateJWT(username string) (string, error) {
	if m.CreateJWTFunc != nil {
		return m.CreateJWTFunc(username)
	}
	return "", nil
}

func (m *MockTokenService) GetUserDataByJWT(cookies []*http.Cookie) (utils.UserData, error) {
	if m.GetUserDataByJWTFunc != nil {
		return m.GetUserDataByJWTFunc(cookies)
	}
	return utils.UserData{}, nil
}

func (m *MockTokenService) GetUserByJWT(cookies []*http.Cookie) (model.User, error) {
	if m.GetUserByJWTFunc != nil {
		return m.GetUserByJWTFunc(cookies)
	}
	return model.User{}, nil
}

func (m *MockTokenService) IsAuthorized(cookies []*http.Cookie) error {
	if m.IsAuthorizedFunc != nil {
		return m.IsAuthorizedFunc(cookies)
	}
	return nil
}

type MockUserRepository struct {
	GetUserByUsernameFunc func(username string) (model.User, error)
	CreateUserFunc        func(username, name, password string) error
}

func (m *MockUserRepository) GetUserByUsername(username string) (model.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return model.User{}, nil
}

func (m *MockUserRepository) CreateUser(username, name, password string) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(username, name, password)
	}
	return nil
}
