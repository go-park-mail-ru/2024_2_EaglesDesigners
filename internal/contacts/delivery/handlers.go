package delivery

import (
	"context"
	"encoding/json"
	"html"
	"net/http"
	"sync"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/repository"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на получение контактов")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	contacts, err := d.usecase.GetContacts(ctx, user.Username)
	if err != nil {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	log.Println("контакты получены")

	var contactsDTO []models.ContactRespDTO

	for _, contact := range contacts {
		contactsDTO = append(contactsDTO, convertContactToDTO(contact))
	}

	response := models.GetContactsRespDTO{
		Contacts: contactsDTO,
	}

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("контакты успешно отправлены")

	responser.SendStruct(ctx, w, response, http.StatusCreated)
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на добавление контакта")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := json.NewDecoder(r.Body).Decode(&contactCreds); err != nil {
		log.Errorf("в теле запросе нет необходимых тегов")
		responser.SendError(ctx, w, invalidJSONError, http.StatusBadRequest)
		return
	}

	if err := validator.Check(contactCreds); err != nil {
		log.Errorf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	var contactData models.ContactData

	contactData.UserID = user.ID.String()
	contactData.ContactUsername = contactCreds.Username

	contact, err := d.usecase.AddContact(ctx, contactData)
	if err != nil {
		if err == repo.ErrContactAlreadyExist {
			responser.SendError(ctx, w, "Contact already exists", http.StatusBadRequest)
			return
		}

		responser.SendError(ctx, w, "Failed to create contact", http.StatusBadRequest)
		return
	}

	response := convertContactToDTO(contact)

	if err := validator.Check(response); err != nil {
		log.Errorf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	log.Println("Contact delivery: контакт успешно создан")

	responser.SendStruct(ctx, w, response, http.StatusCreated)
}

// DeleteContactHandler godoc
// @Summary Delete contact
// @Description Deletes user contact.
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param credentials body models.ContactReqDTO true "Credentials for delete user contact"
// @Success 200 {object} responser.SuccessResponse "Contact deleted"
// @Failure 400 {object} responser.ErrorResponse "Failed to delete contact"
// @Failure 401 {object} responser.ErrorResponse "Unauthorized"
// @Router /contacts [delete]
func (d *Delivery) DeleteContactHandler(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на удаление контакта")

	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := json.NewDecoder(r.Body).Decode(&contactCreds); err != nil {
		log.Errorf("в теле запросе нет необходимых тегов")
		responser.SendError(ctx, w, invalidJSONError, http.StatusBadRequest)
		return
	}

	var contactData models.ContactData

	contactData.UserID = user.ID.String()
	contactData.ContactUsername = contactCreds.Username

	err := d.usecase.DeleteContact(ctx, contactData)
	if err != nil {
		responser.SendError(ctx, w, "Failed to delete contact", http.StatusBadRequest)
		return
	}

	log.Println("контакт успешно удален")

	responser.SendOK(w, "contact deleted", http.StatusOK)
}

func convertContactToDTO(contact models.Contact) models.ContactRespDTO {
	var safeName *string
	if contact.Name != nil {
		safeName = new(string)
		*safeName = html.EscapeString(*contact.Name)
	}

	var safeAvatarURL *string
	if contact.AvatarURL != nil {
		safeAvatarURL = new(string)
		*safeAvatarURL = html.EscapeString(*contact.AvatarURL)
	}

	return models.ContactRespDTO{
		ID:        contact.ID,
		Username:  html.EscapeString(contact.Username),
		Name:      safeName,
		AvatarURL: safeAvatarURL,
	}
}
