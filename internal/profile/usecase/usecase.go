package usecase

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/multipartHelper"
)

type Repository interface {
	GetProfileByUsername(ctx context.Context, username string) (models.ProfileDataDAO, error)
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
	profile := convertProfileFromDTO(profileDTO)

	avatarURL, err := u.repo.UpdateProfile(ctx, profile)
	if err != nil {
		log.Printf("Не удалось обновить профиль: %v", err)
		return err
	}

	if profile.Avatar != nil {
		err := multipartHepler.RewritePhoto(*profile.Avatar, *avatarURL)
		if err != nil {
			log.Printf("Не удалось перезаписать аватарку: %v", err)
			return err
		}
	}

	return nil
}

func (u *Usecase) GetProfile(ctx context.Context, username string) (models.ProfileData, error) {
	profileDataDAO, err := u.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		log.Printf("Не удалось получить профиль: %v", err)
		return models.ProfileData{}, err
	}

	log.Println("Profile usecase: данные получены")

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
		Username:  dto.Username,
		Name:      dto.Name,
		Bio:       dto.Bio,
		Avatar:    dto.Avatar,
		Birthdate: dto.Birthdate,
	}
}
