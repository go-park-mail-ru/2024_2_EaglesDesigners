package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/custom_error"
	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/validator"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	invalidJSONError  = "Invalid format JSON"
	responseError     = "Failed to create response"
	userNotFoundError = "User not found"
)

var noPerm error = &customerror.NoPermissionError{User: "Alice", Area: "секретная зона"}

type ChatDelivery struct {
	service chatlist.ChatUsecase
}

func NewChatDelivery(service chatlist.ChatUsecase) *ChatDelivery {
	return &ChatDelivery{
		service: service,
	}
}

func init() {
	prometheus.MustRegister(requestGetUserChatsHandlerDuration, requestAddNewChatDuration, requestAddUsersIntoChatDuration,
		requestDeleteUserFromChatDuration, requestDeleteUserFromChatDuration, requestLeaveChatDuration, requestDeleteChatOrGroupDuration,
		requestUpdateGroupDuration, requestGetChatInfoDuration, requestAddBranchDuration, requestSearchChatsDuration)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики для чатов зарегистрированы")
}

var requestGetUserChatsHandlerDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "GetUserChatsHandler_request_duration_seconds",
		Help: "/chats",
	},
	[]string{"method"},
)

// GetUserChatsHandler выдает чаты пользователя в query указать страницу ?page=
//
// GetUserChatsHandler godoc
// @Summary Get chats of user
// @Tags chat
// @Produce json
// @Success 200 {object} model.ChatsDTO
// @Failure 500	"Не удалось получить сообщения"
// @Router /chats [get]
func (c *ChatDelivery) GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestGetUserChatsHandlerDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Пришёл запрос на получения чатов с параметрами: %v", r.URL.Query())

	chats, err := c.service.GetChats(r.Context(), r.Cookies())

	if err != nil {
		fmt.Println(err)

		//вернуть 401
		w.WriteHeader(http.StatusInternalServerError)

		log.Printf("НЕ УДАЛОСЬ ПОЛУЧИТЬ ЧАТЫ. ОШИБКА: %s", err)
		return
	}

	chatsDTO := model.ChatsDTO{
		Chats: chats,
	}

	if err := validator.Check(chatsDTO); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(r.Context(), w, "Invalid data", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(chatsDTO)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

var requestAddNewChatDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "AddNewChat_request_duration_seconds",
		Help: "/addchat",
	},
	[]string{"method"},
)

// AddNewChat godoc
// @Summary Add new chat
// @Tags chat
// @Accept json
// @Param chat body model.ChatDTOInput true "Chat info"
// @Success 201 {object} model.ChatDTOOutput "Чат создан"
// @Failure 400 {object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} responser.ErrorResponse "Не удалось добавить чат / группу"
// @Router /addchat [post]
func (c *ChatDelivery) AddNewChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAddNewChatDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()

	var chatDTO model.ChatDTOInput

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Chat delivery: не удалось распарсить запрос: ", err)
		responser.SendError(ctx, w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	jsonString := r.FormValue("chat_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &chatDTO); err != nil {
			responser.SendError(ctx, w, invalidJSONError, http.StatusBadRequest)
			return
		}
	}

	if err := validator.Check(chatDTO); err != nil {
		log.Printf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		responser.SendError(ctx, w, "Failed to get avatar", http.StatusBadRequest)
		return
	}
	defer func() {
		if avatar != nil {
			avatar.Close()
		}
	}()

	if avatar != nil {
		chatDTO.Avatar = &avatar
	}

	returnChat, err := c.service.AddNewChat(r.Context(), r.Cookies(), chatDTO)
	if err != nil {
		log.Printf("Не удалось добавить чат: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось добавить чат: %v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(returnChat); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, returnChat, 201)
}

var requestAddUsersIntoChatDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "AddUsersIntoChat_request_duration_seconds",
		Help: "/chat/{chatId}/addusers",
	},
	[]string{"method"},
)

// AddUsersIntoChat godoc
// @Summary Добавить пользователей в чат
// @Tags chat
// @Accept json
// @Param users body model.AddUsersIntoChatDTO true "Пользователи на добавление"
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Пользователи добавлены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить пользователей"
// @Router /chat/{chatId}/addusers [post]
func (c *ChatDelivery) AddUsersIntoChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAddUsersIntoChatDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> AddUsersIntoChat: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> AddUsersIntoChat: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	var usersToAdd model.AddUsersIntoChatDTO
	err = json.NewDecoder(r.Body).Decode(&usersToAdd)
	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось распарсить Json: %v", err), http.StatusBadRequest)
		return
	}

	if err := validator.Check(usersToAdd); err != nil {
		log.Printf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	addedUsers, err := c.service.AddUsersIntoChatWithCheckPermission(r.Context(), usersToAdd.UsersId, chatUUID)

	if err != nil {
		log.Printf("Не удалось добавить пользователей в чат: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось добавить пользователей в чат: %v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(addedUsers); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, addedUsers, 200)
}

var requestDeleteUsersFromChatDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "DeleteUsersFromChat_request_duration_seconds",
		Help: "/chat/{chatId}/delusers",
	},
	[]string{"method"},
)

// DeleteUsersFromChat godoc
// @Summary Удалить пользователей из чата
// @Tags chat
// @Accept json
// @Param users body model.DeleteUsersFromChatDTO true "Пользователи на добавление"
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} model.DeletdeUsersFromChatDTO "Пользователи удалены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить пользователей"
// @Router /chat/{chatId}/delusers [delete]
func (c *ChatDelivery) DeleteUsersFromChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestDeleteUsersFromChatDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid:", err)

		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}
	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не переданы параметры", http.StatusInternalServerError)
		return
	}

	var usersToDelete model.DeleteUsersFromChatDTO
	err = json.NewDecoder(r.Body).Decode(&usersToDelete)
	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.Check(usersToDelete); err != nil {
		log.Printf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	delUsers, err := c.service.DeleteUsersFromChat(r.Context(), user.ID, chatUUID, usersToDelete)

	if err != nil {
		log.Printf("Не удалось добавить пользователей в чат: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось добавить пользователей в чат: %v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(delUsers); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, delUsers, 200)
}

var requestDeleteUserFromChatDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "DeleteUserFromChat_request_duration_seconds",
		Help: "/chat/{chatId}/deluser/{userId}",
	},
	[]string{"method"},
)

// DeleteUserFromChat godoc
// @Summary Удалить пользователя из чата
// @Tags chat
// @Param userId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} model.DeletdeUsersFromChatDTO "Пользователь удален"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить пользователей"
// @Router /chat/{chatId}/deluser/{userId} [delete]
func (c *ChatDelivery) DeleteUserFromChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestDeleteUserFromChatDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userToDelId := vars["userId"]

	userToDelUUID, err := uuid.Parse(userToDelId)
	if err != nil {
		log.Errorf("не удалось распарсить messageId: %v", err)
		responser.SendError(ctx, w, "invalid messageId", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не переданы параметры", http.StatusInternalServerError)
		return
	}

	delUsers, err := c.service.DeleteUsersFromChat(r.Context(), user.ID, chatUUID, model.DeleteUsersFromChatDTO{
		UsersId: []uuid.UUID{userToDelUUID},
	})

	if err != nil {
		log.Printf("Не удалось добавить пользователей в чат: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось добавить пользователей в чат: %v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(delUsers); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, delUsers, 200)
}

var requestLeaveChatDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "LeaveChat_request_duration_seconds",
		Help: "/chat/{chatId}/leave",
	},
	[]string{"method"},
)

// LeaveChat godoc
// @Summary Выйти из чата
// @Tags chat
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Пользователь вышел из чата"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Запрещено"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить пользователей"
// @Router /chat/{chatId}/leave [delete]
func (c *ChatDelivery) LeaveChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestLeaveChatDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> DeleteUsersFromChat: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не переданы параметры", http.StatusInternalServerError)
		return
	}

	err = c.service.UserLeaveChat(ctx, user.ID, chatUUID)
	if err != nil {
		if errors.As(err, &noPerm) {
			w.WriteHeader(http.StatusForbidden)
			responser.SendError(ctx, w, fmt.Sprintf("Запрещено: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	responser.SendOK(w, "Пользователь вышел из чата", http.StatusOK)
}

var requestDeleteChatOrGroupDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "DeleteChatOrGroup_request_duration_seconds",
		Help: "/chat/{chatId}/delete",
	},
	[]string{"method"},
)

// DeleteChatOrGroup godoc
// @Summary Удаличть чат или группу
// @Tags chat
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} SuccessfullSuccess "Чат удалён"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Failure 500	{object} responser.ErrorResponse "Не удалось удалить чат"
// @Router /chat/{chatId}/delete [delete]
func (c *ChatDelivery) DeleteChatOrGroup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestDeleteChatOrGroupDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> DeleteChatOrGroup: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> DeleteChatOrGroup: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не получены нужные параметры", http.StatusInternalServerError)
		return
	}

	log.Printf("Chat delivery -> DeleteChatOrGroup: пришёл запрос на удаление чата %v от пользователя %v", chatUUID, user.ID)

	err = c.service.DeleteChat(r.Context(), chatUUID, user.ID)

	if err != nil {
		if errors.As(err, &noPerm) {
			w.WriteHeader(http.StatusForbidden)
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	sendSuccess(ctx, w)
}

func getChatIdFromContext(ctx context.Context) (uuid.UUID, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	mapVars, ok := ctx.Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		return uuid.UUID{}, errors.New("Не удалось достать переменные из контекста")
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting getting all messages for chat: %v", chatUUID)

	if err != nil {
		return uuid.UUID{}, err
	}

	return chatUUID, nil
}

var requestUpdateGroupDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "UpdateGroup_request_duration_seconds",
		Help: "/chat/{chatId} ",
	},
	[]string{"method"},
)

// UpdateGroup godoc
// @Summary Обновляем фото и имя
// @Description Update bio, avatar, name or birthdate of user.
// @Tags chat
// @Accept multipart/form-data
// @Security BearerAuth
// @Param chat_data body model.ChatUpdate true "JSON representation of chat data"
// @Param avatar formData file false "group avatar" example:"/2024_2_eaglesDesigners/uploads/chat/f0364477-bfd4-496d-b639-d825b009d509.png"
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} model.ChatUpdateOutput "Чат обновлен"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Failure 500	{object} responser.ErrorResponse "Не удалось обновчить чат"
// @Router /chat/{chatId} [put]
func (c *ChatDelivery) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestUpdateGroupDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> UpdateGroup: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> UpdateGroup: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не получены нужные параметры", http.StatusInternalServerError)
		return
	}

	var chatUpdate model.ChatUpdate

	jsonString := r.FormValue("chat_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &chatUpdate); err != nil {
			responser.SendError(ctx, w, invalidJSONError, http.StatusBadRequest)
			return
		}
	}

	if err := validator.Check(chatUpdate); err != nil {
		log.Printf("входные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		responser.SendError(ctx, w, "Failed to get avatar", http.StatusBadRequest)
		return
	}
	defer func() {
		if avatar != nil {
			avatar.Close()
		}
	}()

	if avatar != nil {
		log.Println("Chat delivery -> UpdateGroup: обновление аватарки")
		chatUpdate.Avatar = &avatar
	}

	log.Printf("Chat delivery -> UpdateGroup: пришёл запрос на изменение чата %v от пользователя %v", chatUUID, user.ID)

	updatedChat, err := c.service.UpdateChat(r.Context(), chatUUID, chatUpdate, user.ID)

	if err != nil {
		if errors.As(err, &noPerm) {
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("Внутренняя ошибка: %v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(updatedChat); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, updatedChat, 200)
}

func sendSuccess(ctx context.Context, w http.ResponseWriter) {
	responser.SendStruct(ctx, w, SuccessfullSuccess{Success: "Произошёл успешный успех"}, 200)
}

type SuccessfullSuccess struct {
	Success string `json:success`
}

var requestGetChatInfoDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "GetChatInfo_request_duration_seconds",
		Help: "/chat/{chatId} ",
	},
	[]string{"method"},
)

// GetChatInfo godoc
// @Summary Получаем пользователей и последние сообщении чата
// @Tags chat
// @Security BearerAuth
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")// @Success 200 {object} model.ChatUpdateOutput "Чат обновлен"
// @Success 200 {object} model.ChatInfoDTO "Пользователи чата"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Failure 500	{object} responser.ErrorResponse "Не удалось получить учатсников"
// @Router /chat/{chatId} [get]
func (c *ChatDelivery) GetChatInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestGetChatInfoDuration, r.Method)
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> GetUsersFromChat: error parsing chat uuid:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Chat delivery -> GetUsersFromChat: error parsing chat uuid: %v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, "Не получены нужные параметры", http.StatusInternalServerError)
		return
	}

	users, err := c.service.GetChatInfo(ctx, chatUUID, user.ID)
	if err != nil {
		if errors.Is(err, noPerm) {
			responser.SendError(ctx, w, err.Error(), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(users); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, users, http.StatusOK)
}

var requestAddBranchDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "AddBranch_request_duration_seconds",
		Help: "/chat/{chatId}/{messageId}/branch",
	},
	[]string{"method"},
)

// UpdateGroup godoc
// @Summary Добавить ветку к сообщению в чате
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} model.AddBranch "Ветка добавлена"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Router /chat/{chatId}/{messageId}/branch [post]
func (c *ChatDelivery) AddBranch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAddBranchDuration, r.Method)
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("пришел запрос /chat/{chatid}/{messageId}/branch [post]")

	vars := mux.Vars(r)
	chatID := vars["chatId"]
	messageID := vars["messageId"]

	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		log.Errorf("не удалось распарсить messageId: %v", err)
		responser.SendError(ctx, w, "invalid messageId", http.StatusBadRequest)
	}

	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		log.Errorf("не удалось распарсить chatid: %v", err)
		responser.SendError(ctx, w, "invalid messageId", http.StatusBadRequest)
	}

	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		responser.SendError(ctx, w, userNotFoundError, http.StatusNotFound)
		return
	}

	branch, err := c.service.AddBranch(ctx, chatUUID, messageUUID, user.ID)
	if err != nil {
		log.Errorf("не удалось создать ветку: %v", err)
		responser.SendError(ctx, w, "invalid data", http.StatusBadRequest)
	}

	responser.SendStruct(ctx, w, branch, http.StatusCreated)
}

var requestSearchChatsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "SearchChats_request_duration_seconds",
		Help: "/chat/search",
	},
	[]string{"method"},
)

// SearchChats ищет чаты по названию, в query указать ключевое слово ?key_word=
//
// SearchChats godoc
// @Summary Поиск чатов пользователя и глобальных каналов по названию
// @Tags chat
// @Produce json
// @Security BearerAuth
// @Param key_word query string false "Ключевое слово для поиска"
// @Success 200 {object} model.SearchChatsDTO
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} responser.ErrorResponse "Нет полномочий"
// @Failure 500	"Не удалось получить сообщения"
// @Router /chat/search [get]
func (c *ChatDelivery) SearchChats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestSearchChatsDuration, r.Method)
	}()

	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Debugf("Пришёл запрос на получения чатов с параметрами: %v", r.URL.Query())
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

	output, err := c.service.SearchChats(ctx, user.ID, keyWord)
	if err != nil {
		log.Errorf("НЕ УДАЛОСЬ ПОЛУЧИТЬ ЧАТЫ. ОШИБКА: %v", err)
		responser.SendError(ctx, w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := validator.Check(output); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	responser.SendStruct(ctx, w, output, http.StatusOK)
}
