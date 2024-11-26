package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/usecase"
	mock_repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.octolab.org/pointer"
)

var errRepo = errors.New("repo error")

func TestUpdateProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repo.NewMockRepository(ctrl)

	usecase := usecase.New(mockRepo)

	tests := []struct {
		name          string
		profileDTO    models.UpdateProfileRequestDTO
		prepareMock   func()
		expectedError error
	}{
		{
			name: "успешное обновление профиля",
			profileDTO: models.UpdateProfileRequestDTO{
				ID:   uuid.New(),
				Name: pointer.ToString("Test User"),
			},
			prepareMock: func() {
				mockRepo.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					Return(nil, nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "ошибка при обновлении профиля в репозитории",
			profileDTO: models.UpdateProfileRequestDTO{
				ID:   uuid.New(),
				Name: pointer.ToString("Test User"),
			},
			prepareMock: func() {
				mockRepo.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					Return(nil, nil, errRepo)
			},
			expectedError: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()

			err := usecase.UpdateProfile(context.Background(), tt.profileDTO)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: '%v', got: '%v'", tt.expectedError, err)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repo.NewMockRepository(ctrl)

	usecase := usecase.New(mockRepo)

	tests := []struct {
		name            string
		id              uuid.UUID
		prepareMock     func()
		expectedError   error
		expectedProfile models.ProfileData
	}{
		{
			name: "success",
			id:   uuid.New(),
			prepareMock: func() {
				mockRepo.EXPECT().
					GetProfileByUsername(gomock.Any(), gomock.Any()).
					Return(models.ProfileDataDAO{
						Name:      pointer.ToString("test"),
						Bio:       pointer.ToString("test"),
						Birthdate: &time.Time{},
					}, nil)
			},
			expectedError: nil,
			expectedProfile: models.ProfileData{
				Name:      pointer.ToString("test"),
				Bio:       pointer.ToString("test"),
				Birthdate: &time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()

			_, err := usecase.GetProfile(context.Background(), tt.id)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: '%v', got: '%v'", tt.expectedError, err)
			}
		})
	}
}
