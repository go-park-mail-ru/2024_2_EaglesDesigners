package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/multipartHelper"
	"github.com/google/uuid"
)

type Repository interface {
	GetProfileByUsername(ctx context.Context, id uuid.UUID) (models.ProfileDataDAO, error)
	UpdateProfile(ctx context.Context, profile models.Profile) (avatarURL *string, err error)
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) UpdateProfile(ctx context.Context, profileDTO models.UpdateProfileRequestDTO) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	profile := convertProfileFromDTO(profileDTO)

	avatarURL, err := u.repo.UpdateProfile(ctx, profile)
	if err != nil {
		log.Errorf("не удалось обновить профиль: %v", err)
		return err
	}

	if profile.Avatar != nil {
		err := multipartHepler.RewritePhoto(*profile.Avatar, *avatarURL)
		if err != nil {
			log.Errorf("не удалось перезаписать аватарку: %v", err)
			return err
		}
	}

	return nil
}

func (u *Usecase) GetProfile(ctx context.Context, id uuid.UUID) (models.ProfileData, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	profileDataDAO, err := u.repo.GetProfileByUsername(ctx, id)
	if err != nil {
		log.Errorf("Не удалось получить профиль: %v", err)
		return models.ProfileData{}, err
	}

	log.Println("данные получены")

	profileData := convertProfileDataFromDAO(profileDataDAO)

	return profileData, nil
}

func convertProfileDataFromDAO(dao models.ProfileDataDAO) models.ProfileData {
	return models.ProfileData{
		Name:       dao.Name,
		Bio:        dao.Bio,
		Birthdate:  dao.Birthdate,
		AvatarPath: dao.AvatarPath,
	}
}

func convertProfileFromDTO(dto models.UpdateProfileRequestDTO) models.Profile {
	return models.Profile{
		ID:        dto.ID,
		Name:      dto.Name,
		Bio:       dto.Bio,
		Avatar:    dto.Avatar,
		Birthdate: dto.Birthdate,
	}
}
