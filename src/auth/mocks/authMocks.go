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
	return m.AuthenticateFunc(username, password)
}

func (m *MockAuthService) Registration(username, name, password string) error {
	return m.RegistrationFunc(username, name, password)
}

func (m *MockAuthService) GetUserDataByUsername(username string) (utils.UserData, error) {
	return m.GetUserDataFunc(username)
}

type MockTokenService struct {
	CreateJWTFunc        func(username string) (string, error)
	GetUserDataByJWTFunc func(cookies []*http.Cookie) (utils.UserData, error)
	GetUserByJWTFunc     func(cookies []*http.Cookie) (model.User, error)
	IsAuthorizedFunc     func(cookies []*http.Cookie) error
}

func (m *MockTokenService) CreateJWT(username string) (string, error) {
	return m.CreateJWTFunc(username)
}

func (m *MockTokenService) GetUserDataByJWT(cookies []*http.Cookie) (utils.UserData, error) {
	return m.GetUserDataByJWTFunc(cookies)
}

func (m *MockTokenService) GetUserByJWT(cookies []*http.Cookie) (model.User, error) {
	return m.GetUserByJWTFunc(cookies)
}

func (m *MockTokenService) IsAuthorized(cookies []*http.Cookie) error {
	return m.IsAuthorizedFunc(cookies)
}

type MockUserRepository struct {
	GetUserByUsernameFunc func(username string) (model.User, error)
	CreateUserFunc        func(username, name, password string) error
}

func (m *MockUserRepository) GetUserByUsername(username string) (model.User, error) {
	return m.GetUserByUsernameFunc(username)
}

func (m *MockUserRepository) CreateUser(username, name, password string) error {
	return m.CreateUserFunc(username, name, password)
}
