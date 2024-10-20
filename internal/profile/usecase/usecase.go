package usecase

import (
	"context"
	"log"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"
)

const (
	defaultAvatarBase64Value = ""
)

type Repository interface {
	GetProfileByUsername(ctx context.Context, username string) (models.ProfileDataDAO, error)
	UpdateProfile(ctx context.Context, profile models.Profile) (avatarURL string, err error)
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
	var avatarChanged bool

	if profileDTO.AvatarBase64 != defaultAvatarBase64Value {
		avatarChanged = true
	}

	profile := convertProfileFromDTO(profileDTO)

	avatarURL, err := u.repo.UpdateProfile(ctx, profile)
	if err != nil {
		log.Printf("Не удалось обновить профиль: %v", err)
		return err
	}

	if avatarChanged {
		err := base64helper.RewritePhoto(profileDTO.AvatarBase64, avatarURL)
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

	profileData := convertProfileDataFromDAO(profileDataDAO)

	return profileData, nil
}

func convertProfileDataFromDAO(dao models.ProfileDataDAO) models.ProfileData {
	var bio *string

	if dao.Bio.Valid {
		bio = &dao.Bio.String
	} else {
		bio = nil
	}

	var avatarBase64 *string

	if dao.AvatarURL.Valid {
		avatarBase64 = &dao.AvatarURL.String
	} else {
		avatarBase64 = nil
	}

	var birthdate *time.Time

	if dao.Birthdate.Valid {
		birthdate = &dao.Birthdate.Time
	} else {
		birthdate = nil
	}

	return models.ProfileData{
		Bio:       bio,
		AvatarURL: avatarBase64,
		Birthdate: birthdate,
	}
}

func convertProfileFromDTO(dto models.UpdateProfileRequestDTO) models.Profile {
	return models.Profile{
		Username:     dto.Username,
		Name:         dto.Name,
		Bio:          dto.Bio,
		AvatarBase64: dto.AvatarBase64,
		Birthdate:    dto.Birthdate,
	}
}
