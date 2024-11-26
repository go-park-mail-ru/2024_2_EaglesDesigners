package repository_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/repository"
	mocks "github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/pgxpool_mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// MockRow реализует интерфейс pgx.Row для подмены результата.
type MockRow struct {
	user models.UserDAO
	err  error
}

func (mr *MockRow) Scan(dest ...interface{}) error {
	if mr.err != nil {
		return mr.err
	}

	// 	0 &user.ID,
	// 	1 &user.Username,
	// 	2 &user.Password,
	// 	3 &user.Version,
	// 	4 &user.Name,
	// 	5 &user.AvatarURL,

	if v, ok := dest[0].(*uuid.UUID); ok {
		*v = mr.user.ID
	}
	if v, ok := dest[1].(*string); ok {
		*v = mr.user.Username
	}
	if v, ok := dest[2].(*string); ok {
		*v = mr.user.Password
	}
	if v, ok := dest[3].(*int64); ok {
		*v = mr.user.Version
	}
	if v, ok := dest[4].(*string); ok {
		*v = mr.user.Name
	}
	if v, ok := dest[5].(*string); ok {
		if mr.user.AvatarURL != nil {
			*v = *mr.user.AvatarURL
		} else {
			*v = "" // Если AvatarURL nil, то присваиваем пустую строку.
		}
	}
	return nil
}

func TestGetUserByUsername_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockPgxIface(ctrl)
	repo := repository.NewRepository(mockDB)

	ctx := context.Background()
	username := "testuser"
	expectedUser := models.UserDAO{
		ID:       uuid.New(), // Генерация нового UUID
		Username: username,
		Name:     "test",
		Password: "hashed_password",
		Version:  1,
	}

	// Настройка мока для успешного выполнения запроса
	mockDB.EXPECT().
		QueryRow(ctx, gomock.Any(), username).
		Return(&MockRow{user: expectedUser, err: nil})

	user, err := repo.GetUserByUsername(ctx, username)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user != expectedUser {
		t.Errorf("expected user %+v, got %+v", expectedUser, user)
	}
}

func TestGetUserByUsername_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPgxIface(ctrl)
	repo := repository.NewRepository(mockDB)

	ctx := context.Background()
	username := "nonexistentuser"

	// Настройка мока для возвращения ошибки
	mockDB.EXPECT().
		QueryRow(ctx, gomock.Any(), username).
		Return(&MockRow{err: pgx.ErrNoRows})

	user, err := repo.GetUserByUsername(ctx, username)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if user != (models.UserDAO{}) {
		t.Errorf("expected empty user, got %+v", user)
	}
}

func TestCreateUser_Success(t *testing.T) {

}
