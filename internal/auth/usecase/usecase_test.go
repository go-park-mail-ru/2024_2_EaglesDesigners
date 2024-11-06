package usecase_test

import (
	"errors"
	"testing"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	mocks "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := &mocks.Mockrepository{
		GetUserByUsernameFunc: func(username string) (models.User, error) {
			return repo.User{
				Username: username,
				Password: usecase.HashPassword("pass1"),
			}, nil
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)

	assert.True(t, authUC.Authenticate("user11", "pass1"))
}

func TestAuthenticate_Failure_UserNotFound(t *testing.T) {
	mockRepo := &mocks.Mockrepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{}, errors.New("user does not exist")
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)

	assert.False(t, authUC.Authenticate("unknown_user", "pass1"))
}

func TestRegistration_Success(t *testing.T) {
	mockRepo := &mocks.Mockrepository{
		CreateUserFunc: func(username, name, password string) error {
			return nil
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)
	err := authUC.Registration("new_user", "John Doe", "password1")

	assert.NoError(t, err)
}

func TestRegistration_Failure_UserExists(t *testing.T) {
	mockRepo := &mocks.Mockrepository{
		CreateUser: func(username, name, password string) error {
			return errors.New("user does not exist")
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)
	err := authUC.Registration("existing_user", "John Doe", "pass1")

	assert.Error(t, err)
}

func TestGetUserDataByUsername_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				ID:       1,
				Username: username,
				Name:     "John Doe",
			}, nil
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)
	userData, err := authUC.GetUserDataByUsername("user1")

	assert.NoError(t, err)
	assert.Equal(t, "user1", userData.Username)
	assert.Equal(t, "John Doe", userData.Name)
}

func TestGetUserDataByUsername_Failure_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{}, errors.New("user does not exist")
		},
	}
	authUC := usecase.NewUsecase(mockRepo, nil)
	userData, err := authUC.GetUserDataByUsername("unknown_user")

	assert.Error(t, err)
	assert.Equal(t, usecase.UserData{}, userData)
}
