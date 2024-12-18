package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/models"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/multipartHelper"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mocks.go

type Repository interface {
	GetProfileByUsername(ctx context.Context, id uuid.UUID) (models.ProfileDataDAO, error)
	UpdateProfile(ctx context.Context, profile models.Profile) (avatarNewURL *string, avatarOldURL *string, err error)
}

type FilesUsecase interface {
	RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) error
	DeletePhoto(ctx context.Context, fileIDStr string) error
	IsImage(header multipart.FileHeader) error
}

type Usecase struct {
	filesUC FilesUsecase
	repo    Repository
}

func New(filesUC FilesUsecase, repo Repository) *Usecase {
	return &Usecase{
		filesUC: filesUC,
		repo:    repo,
	}
}

func (u *Usecase) UpdateProfile(ctx context.Context, profileDTO models.UpdateProfileRequestDTO) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	profile := convertProfileFromDTO(profileDTO)

	if profile.Avatar != nil {
		if profile.AvatarHeader == nil {
			log.Errorf("нет хедера")
			return errors.New("нет хедера")
		}
		if err := u.filesUC.IsImage(*profile.AvatarHeader); err != nil {
			log.WithError(err).Errorf("аватарка не картинка")
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
		if err := u.filesUC.RewritePhoto(ctx, *profile.Avatar, *profile.AvatarHeader, *avatarNewURL); err != nil {
			log.Errorf("не удалось сохранить аватарку: %v", err)
			return err
		}
		log.Printf("Аватарка успешно сохранена")

		if avatarOldURL != nil {
			go func() {
				log.Printf("Удаление старой аватарки %s", *avatarOldURL)
				if err := u.filesUC.DeletePhoto(ctx, *avatarOldURL); err != nil {
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
		AvatarHeader: dto.AvatarHeader,
		DeleteAvatar: dto.DeleteAvatar,
		Birthdate:    dto.Birthdate,
	}
}
