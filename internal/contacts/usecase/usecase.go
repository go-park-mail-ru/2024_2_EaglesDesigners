package usecase

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"
	"github.com/google/uuid"
)

type Repository interface {
	GetContacts(ctx context.Context, username string) (contacts []models.UserDAO, err error)
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) GetContacts(ctx context.Context, username string) (contacts []models.User, err error) {
	log.Println("UC before get con")
	contactsDAO, err := u.repo.GetContacts(ctx, username)
	if err != nil {
		log.Printf("Не удалось получить контакты: %v", err)
		return contacts, err
	}

	log.Println("UC after get con")
	log.Println("uc contacts dao", contactsDAO)

	for _, contactDAO := range contactsDAO {
		log.Println("uc in for", contactDAO.Username, contactDAO.Name)
		contact, err := convertUserFromDAO(contactDAO)
		if err != nil {
			return contacts, err
		}
		log.Println("uc ")

		contacts = append(contacts, *contact)
	}

	log.Println("US done")

	return contacts, nil
}

// func (u *Usecase) UpdateProfile(ctx context.Context, profileDTO models.UpdateProfileRequestDTO) error {
// var avatarChanged bool

// if profileDTO.AvatarBase64 != nil {
// 	avatarChanged = true
// }

// profile := convertProfileFromDTO(profileDTO)

// avatarURL, err := u.repo.UpdateProfile(ctx, profile)
// if err != nil {
// 	log.Printf("Не удалось обновить профиль: %v", err)
// 	return err
// }

// if avatarChanged {
// 	err := base64helper.RewritePhoto(*profileDTO.AvatarBase64, *avatarURL)
// 	if err != nil {
// 		log.Printf("Не удалось перезаписать аватарку: %v", err)
// 		return err
// 	}
// }

// return nil
// }

func convertUserFromDAO(dao models.UserDAO) (*models.User, error) {
	avatarUUID, err := uuid.Parse(*dao.AvatarURL)
	if err != nil {
		return &models.User{}, err
	}

	avatarBase64, err := base64helper.ReadPhotoBase64(avatarUUID)
	if err != nil {
		return &models.User{}, err
	}

	user := models.User{
		Username:     dao.Username,
		Name:         dao.Name,
		AvatarBase64: &avatarBase64,
	}

	return &user, nil
}
