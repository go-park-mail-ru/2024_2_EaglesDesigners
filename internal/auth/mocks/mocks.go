package mocks

import (
	"net/http"

	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	authUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	jwtUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type MockUsecase struct {
	AuthenticateFunc func(username, password string) bool
	RegistrationFunc func(username, name, password string) error
	GetUserDataFunc  func(username string) (authUC.UserData, error)
}

func (m *MockUsecase) Authenticate(username, password string) bool {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(username, password)
	}
	return false
}

func (m *MockUsecase) Registration(username, name, password string) error {
	if m.RegistrationFunc != nil {
		return m.RegistrationFunc(username, name, password)
	}
	return nil
}

func (m *MockUsecase) GetUserDataByUsername(username string) (authUC.UserData, error) {
	if m.GetUserDataFunc != nil {
		return m.GetUserDataFunc(username)
	}
	return authUC.UserData{}, nil
}

type MockRepository struct {
	GetUserByUsernameFunc func(username string) (repo.User, error)
	CreateUserFunc        func(username, name, password string) error
}

func (m *MockRepository) GetUserByUsername(username string) (repo.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return repo.User{}, nil
}

func (m *MockRepository) CreateUser(username, name, password string) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(username, name, password)
	}
	return nil
}

type MockTokenUsecase struct {
	CreateJWTFunc        func(username string) (string, error)
	GetUserDataByJWTFunc func(cookies []*http.Cookie) (jwtUC.UserData, error)
	GetUserByJWTFunc     func(cookies []*http.Cookie) (repo.User, error)
	IsAuthorizedFunc     func(cookies []*http.Cookie) error
}

func (m *MockTokenUsecase) CreateJWT(username string) (string, error) {
	if m.CreateJWTFunc != nil {
		return m.CreateJWTFunc(username)
	}
	return "", nil
}

func (m *MockTokenUsecase) GetUserDataByJWT(cookies []*http.Cookie) (jwtUC.UserData, error) {
	if m.GetUserDataByJWTFunc != nil {
		return m.GetUserDataByJWTFunc(cookies)
	}
	return jwtUC.UserData{}, nil
}

func (m *MockTokenUsecase) GetUserByJWT(cookies []*http.Cookie) (jwtUC.User, error) {
	if m.GetUserByJWTFunc != nil {
		repoUser, err := m.GetUserByJWTFunc(cookies)
		if err != nil {
			return jwtUC.User{}, err
		}

		user := convertToUser(repoUser)
		return user, nil
	}
	return jwtUC.User{}, nil
}

func (m *MockTokenUsecase) IsAuthorized(cookies []*http.Cookie) error {
	if m.IsAuthorizedFunc != nil {
		return m.IsAuthorizedFunc(cookies)
	}
	return nil
}

func convertToUser(u repo.User) jwtUC.User {
	return jwtUC.User{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Version:  u.Version,
	}
}
