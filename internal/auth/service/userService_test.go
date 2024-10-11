package service_test

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/service"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		GetUserByUsernameFunc: func(username string) (model.User, error) {
			return model.User{
				Username: username,
				Password: utils.HashPassword("pass1"), // Пример правильного хеша
			}, nil
		},
	}
	authService := service.NewAuthService(mockRepo, nil)

	assert.True(t, authService.Authenticate("user1", "pass1"))
}

func TestAuthenticate_Failure_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		GetUserByUsernameFunc: func(username string) (model.User, error) {
			return model.User{}, errors.New("user does not exist")
		},
	}
	authService := service.NewAuthService(mockRepo, nil)

	assert.False(t, authService.Authenticate("unknown_user", "pass1"))
}

func TestRegistration_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		CreateUserFunc: func(username, name, password string) error {
			return nil
		},
	}
	authService := service.NewAuthService(mockRepo, nil)
	err := authService.Registration("new_user", "John Doe", "password1")

	assert.NoError(t, err)
}

func TestRegistration_Failure_UserExists(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		CreateUserFunc: func(username, name, password string) error {
			return errors.New("user does not exist")
		},
	}
	authService := service.NewAuthService(mockRepo, nil)
	err := authService.Registration("existing_user", "John Doe", "pass1")

	assert.Error(t, err)
}

func TestGetUserDataByUsername_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		GetUserByUsernameFunc: func(username string) (model.User, error) {
			return model.User{
				ID:       1,
				Username: username,
				Name:     "John Doe",
			}, nil
		},
	}
	authService := service.NewAuthService(mockRepo, nil)
	userData, err := authService.GetUserDataByUsername("user1")

	assert.NoError(t, err)
	assert.Equal(t, "user1", userData.Username)
	assert.Equal(t, "John Doe", userData.Name)
}

func TestGetUserDataByUsername_Failure_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		GetUserByUsernameFunc: func(username string) (model.User, error) {
			return model.User{}, errors.New("user does not exist")
		},
	}
	authService := service.NewAuthService(mockRepo, nil)
	userData, err := authService.GetUserDataByUsername("unknown_user")

	assert.Error(t, err)
	assert.Equal(t, utils.UserData{}, userData)
}
