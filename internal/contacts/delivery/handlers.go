package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
)

const (
	invalidJSONError  = "Invalid format JSON"
	responseError     = "Failed to create response"
	userNotFoundError = "User not found"
)

type usecase interface {
	GetContacts(ctx context.Context, username string) (contacts []models.Contact, err error)
	AddContact(ctx context.Context, contactData models.ContactData) (models.Contact, error)
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

	log.Println("Пришел запрос на получение контактов")

	log.Println("Получение пользователя из jwt")
	user, err := d.token.GetUserDataByJWT(r.Cookies())
	if err != nil {
		log.Println("Пользователь не найден")
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	contacts, err := d.usecase.GetContacts(ctx, user.Username)
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	log.Println("Контакты получены")

	var contactsDTO []models.ContactDTO

	for _, contact := range contacts {
		contactsDTO = append(contactsDTO, convertContactToDTO(contact))
	}

	response := models.GetContactsRespDTO{
		Contacts: contactsDTO,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// AddContactHandler godoc
// @Summary Add new contact
// @Description Create a new contact for the user.
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param credentials body models.AddContactReqDTO true "Credentials for create a new contact"
// @Success 201 {object} models.ContactDTO "Contact created"
// @Failure 400 {object} responser.ErrorResponse "Failed to create contact"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Failure 404 {object} responser.ErrorResponse "User not found"
// @Router /contacts [post]
func (d *Delivery) AddContactHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()

	log.Println("Пришел запрос на добавление контакта")

	userData, err := d.token.GetUserDataByJWT(r.Cookies())
	if err != nil {
		responser.SendErrorResponse(w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.AddContactReqDTO

	if err := json.NewDecoder(r.Body).Decode(&contactCreds); err != nil {
		log.Println("В теле запросе нет необходимых тегов")
		responser.SendErrorResponse(w, invalidJSONError, http.StatusBadRequest)
		return
	}

	var contactData models.ContactData

	contactData.UserID = userData.ID.String()
	contactData.ContactUsername = contactCreds.Username

	contact, err := d.usecase.AddContact(ctx, contactData)
	if err != nil {
		responser.SendErrorResponse(w, "Failed to create contact", http.StatusBadRequest)
		return
	}

	log.Println("Контакт создан")

	response := convertContactToDTO(contact)

	jsonResp, err := json.Marshal(response)
	if err != nil {
		responser.SendErrorResponse(w, responseError, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func convertContactToDTO(contact models.Contact) models.ContactDTO {
	return models.ContactDTO{
		ID:           contact.ID,
		Username:     contact.Username,
		Name:         contact.Name,
		AvatarBase64: contact.AvatarBase64,
	}
}
