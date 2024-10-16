package usecase_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/mocks"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWT_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				ID:       1,
				Username: username,
				Name:     "John Doe",
				Version:  1,
			}, nil
		},
	}
	tokenUC := usecase.NewUsecase(mockRepo)
	token, err := tokenUC.CreateJWT("user1")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestCreateJWT_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{}, errors.New("user not found")
		},
	}
	tokenUC := usecase.NewUsecase(mockRepo)
	token, err := tokenUC.CreateJWT("unknown_user")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestIsAuthorized_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				Username: "user1",
				Version:  1,
			}, nil
		},
	}
	tokenService := usecase.NewUsecase(mockRepo)
	token, _ := tokenService.CreateJWT("user1")
	cookie := &http.Cookie{Name: "access_token", Value: token}
	cookies := []*http.Cookie{cookie}
	err := tokenService.IsAuthorized(cookies)

	assert.NoError(t, err)
}

func TestIsAuthorized_InvalidToken(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	tokenService := usecase.NewUsecase(mockRepo)
	cookie := &http.Cookie{Name: "access_token", Value: "invalid_token"}
	cookies := []*http.Cookie{cookie}
	err := tokenService.IsAuthorized(cookies)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid token")
}

func TestIsAuthorized_TokenExpired(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				Username: username,
				Version:  1,
			}, nil
		},
	}
	tokenService := usecase.NewUsecase(mockRepo)
	expiredPayload := usecase.Payload{
		Sub:     "user1",
		Name:    "John Doe",
		ID:      1,
		Version: 1,
		Exp:     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	header := usecase.Header{Alg: "HS256", Typ: "JWT"}
	headerJSON, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadJSON, _ := json.Marshal(expiredPayload)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)
	expiredToken, _ := usecase.GeneratorJWT(headerEncoded, payloadEncoded)
	cookie := &http.Cookie{Name: "access_token", Value: expiredToken}
	cookies := []*http.Cookie{cookie}
	err := tokenService.IsAuthorized(cookies)

	assert.Error(t, err)
	assert.EqualError(t, err, "token expired")
}

func TestGetUserByJWT_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				ID:       1,
				Username: username,
				Name:     "John Doe",
			}, nil
		},
	}
	tokenService := usecase.NewUsecase(mockRepo)
	token, _ := tokenService.CreateJWT("user1")
	cookie := &http.Cookie{Name: "access_token", Value: token}
	cookies := []*http.Cookie{cookie}
	user, err := tokenService.GetUserByJWT(cookies)

	assert.NoError(t, err)
	assert.Equal(t, "user1", user.Username)
}

func TestGetUserByJWT_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{}, errors.New("user not found")
		},
	}
	tokenService := usecase.NewUsecase(mockRepo)
	token, _ := tokenService.CreateJWT("some_test_user")
	cookie := &http.Cookie{Name: "access_token", Value: token}
	cookies := []*http.Cookie{cookie}
	_, err := tokenService.GetUserByJWT(cookies)

	assert.Error(t, err)
}

func TestGetUserDataByJWT_Success(t *testing.T) {
	mockRepo := &mocks.MockRepository{
		GetUserByUsernameFunc: func(username string) (repo.User, error) {
			return repo.User{
				ID:       1,
				Username: username,
				Name:     "John Doe",
			}, nil
		},
	}
	tokenService := usecase.NewUsecase(mockRepo)
	token, _ := tokenService.CreateJWT("user1")
	cookie := &http.Cookie{Name: "access_token", Value: token}
	cookies := []*http.Cookie{cookie}
	data, err := tokenService.GetUserDataByJWT(cookies)

	assert.NoError(t, err)
	assert.Equal(t, "user1", data.Username)
	assert.Equal(t, "John Doe", data.Name)
}

func TestGetUserDataByJWT_InvalidToken(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	tokenService := usecase.NewUsecase(mockRepo)
	cookie := &http.Cookie{Name: "access_token", Value: "invalid_token"}
	cookies := []*http.Cookie{cookie}
	_, err := tokenService.GetUserDataByJWT(cookies)

	assert.Error(t, err)
}
