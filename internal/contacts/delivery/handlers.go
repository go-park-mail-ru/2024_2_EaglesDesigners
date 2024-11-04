package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/validator"
)

const (
	invalidJSONError  = "Invalid format JSON"
	userNotFoundError = "User not found"
)

type usecase interface {
	GetContacts(ctx context.Context, username string) (contacts []models.Contact, err error)
	AddContact(ctx context.Context, contactData models.ContactData) (models.Contact, error)
	DeleteContact(ctx context.Context, contactData models.ContactData) error
}

type token interface {
	GetUserDataByJWT(cookies []*http.Cookie) (jwt.UserData, error)
}

type Delivery struct {
	usecase usecase
	token   token
	mu      sync.Mutex
}

func New(usecase usecase, token token) *Delivery {
	return &Delivery{
		usecase: usecase,
		token:   token,
	}
}

// GetContactsHandler godoc
// @Summary Get all contacts
// @Description Get all contacts of user.
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.GetContactsRespDTO "Contacts found"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /contacts [get]
func (d *Delivery) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Contact delivery: пришел запрос на получение контактов")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(w, userNotFoundError, http.StatusNotFound)
		return
	}

	contacts, err := d.usecase.GetContacts(ctx, user.Username)
	if err != nil {
		responser.SendError(w, userNotFoundError, http.StatusNotFound)
		return
	}

	log.Println("Contact delivery: контакты получены")

	var contactsDTO []models.ContactRespDTO

	for _, contact := range contacts {
		contactsDTO = append(contactsDTO, convertContactToDTO(contact))
	}

	response := models.GetContactsRespDTO{
		Contacts: contactsDTO,
	}

	if err := validator.Check(response); err != nil {
		log.Printf("Contact delivery: выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("Contact delivery: контакты успешно отправлены")

	responser.SendStruct(w, response, http.StatusCreated)
}

// AddContactHandler godoc
// @Summary Add new contact
// @Description Create a new contact for the user.
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param credentials body models.ContactReqDTO true "Credentials for create a new contact"
// @Success 201 {object} models.ContactRespDTO "Contact created"
// @Failure 400 {object} responser.ErrorResponse "Failed to create contact"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /contacts [post]
func (d *Delivery) AddContactHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	log.Println("Contact delivery: пришел запрос на добавление контакта")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := json.NewDecoder(r.Body).Decode(&contactCreds); err != nil {
		log.Println("Contact delivery: в теле запросе нет необходимых тегов")
		responser.SendError(w, invalidJSONError, http.StatusBadRequest)
		return
	}

	if err := validator.Check(contactCreds); err != nil {
		log.Printf("Contact delivery: входные данные не прошли проверку валидации: %v", err)
		responser.SendError(w, "Invalid data", http.StatusBadRequest)
		return
	}

	var contactData models.ContactData

	contactData.UserID = user.ID.String()
	contactData.ContactUsername = contactCreds.Username

	contact, err := d.usecase.AddContact(ctx, contactData)
	if err != nil {
		responser.SendError(w, "Failed to create contact", http.StatusBadRequest)
		return
	}

	response := convertContactToDTO(contact)

	if err := validator.Check(response); err != nil {
		log.Printf("Contact delivery: выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("Contact delivery: контакт успешно создан")

	responser.SendStruct(w, response, http.StatusCreated)
}

// DeleteContactHandler godoc
// @Summary Delete contact
// @Description Deletes user contact.
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param credentials body models.ContactReqDTO true "Credentials for delete user contact"
// @Success 200 {object} models.ContactRespDTO "Contact deleted"
// @Failure 400 {object} responser.ErrorResponse "Failed to delete contact"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Router /contacts [delete]
func (d *Delivery) DeleteContactHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	log.Println("Contact delivery: пришел запрос на удаление контакта")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := json.NewDecoder(r.Body).Decode(&contactCreds); err != nil {
		log.Println("Contact delivery: в теле запросе нет необходимых тегов")
		responser.SendError(w, invalidJSONError, http.StatusBadRequest)
		return
	}

	var contactData models.ContactData

	contactData.UserID = user.ID.String()
	contactData.ContactUsername = contactCreds.Username

	err := d.usecase.DeleteContact(ctx, contactData)
	if err != nil {
		responser.SendError(w, "Failed to delete contact", http.StatusBadRequest)
		return
	}

	log.Println("Contact delivery: контакт успешно удален")

	responser.SendOK(w, "contact deleted", http.StatusOK)
}

func convertContactToDTO(contact models.Contact) models.ContactRespDTO {
	return models.ContactRespDTO{
		ID:        contact.ID,
		Username:  contact.Username,
		Name:      contact.Name,
		AvatarURL: contact.AvatarURL,
	}
}
