package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	errGroup "golang.org/x/sync/errgroup"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/models"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mocks.go

func init() {
	prometheus.MustRegister(newContactMetric, deleteContactMetric)
}

type Repository interface {
	GetContacts(ctx context.Context, username string) (contacts []models.ContactDAO, err error)
	AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error)
	DeleteContact(ctx context.Context, contactData models.ContactDataDAO) error
	SearchUserContacts(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error)
	SearchGlobalUsers(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error)
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	contactsDAO, err := u.repo.GetContacts(ctx, username)
	if err != nil {
		log.Errorf("не удалось получить контакты: %v", err)
		return contacts, err
	}
	log.Println("данные получены")

	for _, contactDAO := range contactsDAO {
		contact := convertContactFromDAO(contactDAO)

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

var newContactMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_new_contacts",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

func (u *Usecase) AddContact(ctx context.Context, contactData models.ContactData) (models.Contact, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	contactDataDAO := convertContactDataToDAO(contactData)

	contactDAO, err := u.repo.AddContact(ctx, contactDataDAO)
	if err != nil {
		log.Errorf("не получилось создать контакт: %v", err)
		return models.Contact{}, err
	}

	contact := convertContactFromDAO(contactDAO)

	metric.IncMetric(*newContactMetric)
	return contact, nil
}

var deleteContactMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_deleted_contacts",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

func (u *Usecase) DeleteContact(ctx context.Context, contactData models.ContactData) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	contactDataDAO := convertContactDataToDAO(contactData)

	err := u.repo.DeleteContact(ctx, contactDataDAO)
	if err != nil {
		log.Errorf("не получилось удалить контакт: %v", err)
		return err
	}

	metric.IncMetric(*deleteContactMetric)
	return nil
}

func (u *Usecase) SearchContacts(ctx context.Context, userID uuid.UUID, keyWord string) (models.SearchContactsDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Debugf("пришел запрос на поиск контактов от пользователя: %v", userID)

	var g errGroup.Group

	var userContactsDTO []models.ContactRespDTO
	var globalUsersDTO []models.ContactRespDTO

	g.Go(func() error {
		userContacts, err := u.repo.SearchUserContacts(ctx, userID, keyWord)
		if err != nil {
			return err
		}
		log.Debugln("контакты пользователя получены")

		for _, contact := range userContacts {
			contactDTO := convertContactToDTO(contact)

			userContactsDTO = append(userContactsDTO,
				contactDTO)
		}

		return nil
	})

	g.Go(func() error {
		globalUsers, err := u.repo.SearchGlobalUsers(ctx, userID, keyWord)
		if err != nil {
			return err
		}
		log.Debugln("глобальные пользователи получены")

		for _, user := range globalUsers {
			userDTO := convertContactToDTO(user)

			globalUsersDTO = append(globalUsersDTO,
				userDTO)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return models.SearchContactsDTO{}, err
	}

	outputDTO := models.SearchContactsDTO{
		UserContacts: userContactsDTO,
		GlobalUsers:  globalUsersDTO,
	}

	return outputDTO, nil
}

func convertContactFromDAO(dao models.ContactDAO) models.Contact {
	return models.Contact{
		ID:        dao.ID.String(),
		Username:  dao.Username,
		Name:      dao.Name,
		AvatarURL: dao.AvatarURL,
	}
}

func convertContactToDTO(dao models.ContactDAO) models.ContactRespDTO {
	return models.ContactRespDTO{
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
