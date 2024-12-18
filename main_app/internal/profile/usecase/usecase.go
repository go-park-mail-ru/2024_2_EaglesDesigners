package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/models"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/multipartHelper"
	"github.com/google/uuid"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mocks.go

type Repository interface {
	GetProfileByUsername(ctx context.Context, id uuid.UUID) (models.ProfileDataDAO, error)
	UpdateProfile(ctx context.Context, profile models.Profile) (avatarNewURL *string, avatarOldURL *string, err error)
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

	if profile.Avatar != nil {
		if !multipartHepler.IsImageFile(*profile.Avatar) {
			log.Errorf("аватарка не картинка")
			return multipartHepler.ErrNotImage
		}
	}

	avatarNewURL, avatarOldURL, err := u.repo.UpdateProfile(ctx, profile)
	if err != nil {
		log.Errorf("не удалось обновить профиль: %v", err)
		return err
	}

	if avatarNewURL != nil && profile.Avatar != nil {
		log.Printf("Сохранение аватарки %s", *avatarNewURL)
		if err := multipartHepler.RewritePhoto(*profile.Avatar, *avatarNewURL); err != nil {
			log.Errorf("не удалось сохранить аватарку: %v", err)
			return err
		}
		log.Printf("Аватарка успешно сохранена")

		if avatarOldURL != nil {
			go func() {
				log.Printf("Удаление старой аватарки %s", *avatarOldURL)
				if err := multipartHepler.RemovePhoto(*avatarOldURL); err != nil {
					log.Errorf("Не удалось удалить файл: %v", err)
				}
				log.Println("Удаление прошло успешно")
			}()
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
		ID:           dto.ID,
		Name:         dto.Name,
		Bio:          dto.Bio,
		Avatar:       dto.Avatar,
		DeleteAvatar: dto.DeleteAvatar,
		Birthdate:    dto.Birthdate,
	}
}
