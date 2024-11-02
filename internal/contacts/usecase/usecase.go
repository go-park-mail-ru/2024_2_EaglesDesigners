package usecase

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
)

type Repository interface {
	GetContacts(ctx context.Context, username string) (contacts []models.ContactDAO, err error)
	AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error)
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) GetContacts(ctx context.Context, username string) (contacts []models.Contact, err error) {
	contactsDAO, err := u.repo.GetContacts(ctx, username)
	if err != nil {
		log.Printf("Не удалось получить контакты: %v", err)
		return contacts, err
	}
	log.Println("Usecase: данные получены")

	for _, contactDAO := range contactsDAO {
		contact := convertContactFromDAO(contactDAO)

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (u *Usecase) AddContact(ctx context.Context, contactData models.ContactData) (models.Contact, error) {
	contactDataDAO := convertContactDataToDAO(contactData)

	contactDAO, err := u.repo.AddContact(ctx, contactDataDAO)
	if err != nil {
		log.Println("Usecase: не получилось создать контакт: ", err)
		return models.Contact{}, err
	}

	contact := convertContactFromDAO(contactDAO)

	return contact, nil
}

func convertContactFromDAO(dao models.ContactDAO) models.Contact {
	return models.Contact{
		ID:        dao.ID.String(),
		Username:  dao.Username,
		Name:      dao.Name,
		AvatarURL: dao.AvatarURL,
	}
}

func convertContactDataToDAO(contactData models.ContactData) models.ContactDataDAO {
	return models.ContactDataDAO{
		UserID:          contactData.UserID,
		ContactUsername: contactData.ContactUsername,
	}
}
