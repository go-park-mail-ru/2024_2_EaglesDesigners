package mocks

import (
	"net/http"

	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type MockTokenUsecase struct {
	CreateJWTFunc        func(username string) (string, error)
	GetUserDataByJWTFunc func(cookies []*http.Cookie) (usecase.UserData, error)
	GetUserByJWTFunc     func(cookies []*http.Cookie) (repo.User, error)
	IsAuthorizedFunc     func(cookies []*http.Cookie) error
}

func (m *MockTokenUsecase) CreateJWT(username string) (string, error) {
	if m.CreateJWTFunc != nil {
		return m.CreateJWTFunc(username)
	}
	return "", nil
}

func (m *MockTokenUsecase) GetUserDataByJWT(cookies []*http.Cookie) (usecase.UserData, error) {
	if m.GetUserDataByJWTFunc != nil {
		return m.GetUserDataByJWTFunc(cookies)
	}
	return usecase.UserData{}, nil
}

func (m *MockTokenUsecase) GetUserByJWT(cookies []*http.Cookie) (usecase.User, error) {
	if m.GetUserByJWTFunc != nil {
		repoUser, err := m.GetUserByJWTFunc(cookies)
		if err != nil {
			return usecase.User{}, err
		}

		user := convertToUser(repoUser)
		return user, nil
	}
	return usecase.User{}, nil
}

func (m *MockTokenUsecase) IsAuthorized(cookies []*http.Cookie) error {
	if m.IsAuthorizedFunc != nil {
		return m.IsAuthorizedFunc(cookies)
	}
	return nil
}

func convertToUser(u repo.User) usecase.User {
	return usecase.User{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Version:  u.Version,
	}
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
