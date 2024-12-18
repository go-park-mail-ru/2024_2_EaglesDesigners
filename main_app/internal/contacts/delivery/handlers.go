package delivery

import (
	"context"
	"html"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/models"
	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/validator"
)

const (
	invalidJSONError  = "Invalid format JSON"
	userNotFoundError = "User not found"
)

//go:generate mockgen -source=handlers.go -destination=mocks/mocks.go

type usecase interface {
	GetContacts(ctx context.Context, username string) (contacts []models.Contact, err error)
	AddContact(ctx context.Context, contactData models.ContactData) (models.Contact, error)
	DeleteContact(ctx context.Context, contactData models.ContactData) error
	SearchContacts(ctx context.Context, userID uuid.UUID, keyWord string) (models.SearchContactsDTO, error)
}

type Delivery struct {
	usecase usecase
	mu      sync.Mutex
}

func New(usecase usecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}

func init() {
	prometheus.MustRegister(requestContactDuration)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики для контактов зарегистрированы")
}

var requestContactDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "request_contact_duration_seconds",
	},
	[]string{"method"},
)

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
// @Router /contacts [get].
func (d *Delivery) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestContactDuration, "GetContactsHandler")
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на получение контактов")

	user, ok := ctx.Value(auth.UserKey).(auth.User)
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

	// responser.SendStruct(ctx, w, response, http.StatusOK)
	jsonResp, err := easyjson.Marshal(response)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusOK)
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
// @Router /contacts [post].
func (d *Delivery) AddContactHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestContactDuration, "AddContactHandler")
	}()

	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на добавление контакта")

	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &contactCreds); err != nil {
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

	// responser.SendStruct(ctx, w, response, http.StatusCreated)
	jsonResp, err := easyjson.Marshal(response)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusCreated)
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
// @Router /contacts [delete].
func (d *Delivery) DeleteContactHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestContactDuration, "DeleteContactHandler")
	}()

	d.mu.Lock()
	defer d.mu.Unlock()
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос на удаление контакта")

	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	var contactCreds models.ContactReqDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &contactCreds); err != nil {
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

// SearchChats ищет контакты по имени или нику, в query указать ключевое слово ?key_word=
//
// SearchChats godoc
// @Summary Поиск контактов пользователя и глобальных пользователей по имени или нику
// @Tags contacts
// @Produce json
// @Security BearerAuth
// @Param key_word query string false "Ключевое слово для поиска"
// @Success 200 {object} models.SearchContactsDTO
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Failure 500	"Не удалось получить контакты"
// @Router /contacts/search [get].
func (d *Delivery) SearchContactsHandler(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestContactDuration, "SearchContactsHandler")
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Debugf("Пришёл запрос на поиск контактов с параметрами: %v", r.URL.Query())
	keyWord := r.URL.Query().Get("key_word")

	if keyWord == "" {
		log.Errorf("key_word отсутствует")
		responser.SendError(ctx, w, "key_word not found", http.StatusBadRequest)
		return
	}

	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	output, err := d.usecase.SearchContacts(ctx, user.ID, keyWord)
	if err != nil {
		log.Errorf("не удалось получить контакты: %v", err)
		responser.SendError(ctx, w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := validator.Check(output); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	// responser.SendStruct(ctx, w, output, http.StatusOK)
	jsonResp, err := easyjson.Marshal(output)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusOK)
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
