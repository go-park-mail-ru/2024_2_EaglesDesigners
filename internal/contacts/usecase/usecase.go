package usecase

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
)

type Repository interface {
	GetContacts(ctx context.Context, username string) (contacts []models.ContactDAO, err error)
	AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error)
	DeleteContact(ctx context.Context, contactData models.ContactDataDAO) error
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
		log.Printf("Contact usecase: не удалось получить контакты: %v", err)
		return contacts, err
	}
	log.Println("Contact usecase: данные получены")

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
		log.Println("Contact usecase: не получилось создать контакт: ", err)
		return models.Contact{}, err
	}

	contact := convertContactFromDAO(contactDAO)

	return contact, nil
}

func (u *Usecase) DeleteContact(ctx context.Context, contactData models.ContactData) error {
	contactDataDAO := convertContactDataToDAO(contactData)

	err := u.repo.DeleteContact(ctx, contactDataDAO)
	if err != nil {
		log.Println("Contact usecase: не получилось удалить контакт: ", err)
		return err
	}

	return nil
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
