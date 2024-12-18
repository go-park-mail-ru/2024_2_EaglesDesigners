package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/custom_error"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/validator"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
)

var noPerm error = &customerror.NoPermissionError{User: "Alice", Area: "секретная зона"}

type MessageController struct {
	usecase usecase.MessageUsecase
}

func NewMessageController(usecase usecase.MessageUsecase) MessageController {
	return MessageController{
		usecase: usecase,
	}
}

func init() {
	prometheus.MustRegister(requestMessageDuration)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики для сообщений зарегистрированы")
}

var requestMessageDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "request_messgae_duration_seconds",
	},
	[]string{"method"},
)

// AddNewMessageHandler godoc
// @Summary Add new message
// @Tags message
// @Accept json
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param message body models.MessageInput true "Message info"
// @Success 201 "Сообщение успешно добавлено"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить сообщение"
// @Router /chat/{chatId}/messages [post]
func (h *MessageController) AddNewMessage(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "AddNewMessage")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	if err != nil {
		//conn.400
		log.Println("Delivery: error during parsing json:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Delivery: error during connection upgrade:%v", err), http.StatusBadRequest)
		return
	}

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting adding new message for chat: %v", chatUUID)

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	log.Println(user)
	if !ok {
		log.Println("Message delivery -> AddNewMessage: нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	var messageDTO models.Message
	err = json.NewDecoder(r.Body).Decode(&messageDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось распарсить Json: %v", err), http.StatusBadRequest)
		return
	}

	err = h.usecase.SendMessage(r.Context(), user, chatUUID, messageDTO)

	if err != nil {
		log.Printf("Не удалось добавить сообщение: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось добавить сообщение: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteMessage godoc
// @Summary Delete message
// @Tags message
// @Param messageId path string true "messageId ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Сообщение успешно удалено"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} customerror.NoPermissionError "Нет доступа"
// @Failure 500	{object} responser.ErrorResponse "Не удалось удалить сообщение"
// @Router /messages/{messageId} [delete]
func (h *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "DeleteMessage")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	chatId := mapVars["messageId"]
	log.Printf("messageId: %s", chatId)
	messageUUID, err := uuid.Parse(chatId)
	if err != nil {
		//conn.400
		log.Printf("Получен кривой Id сообщения %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Получен кривой Id сообщения %v", err), http.StatusBadRequest)
		return
	}
	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		log.Println("нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}
	err = h.usecase.DeleteMessage(ctx, user, messageUUID)

	if err != nil {
		if errors.As(err, &noPerm) {
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("внутренняя ошибка: %v", err), http.StatusInternalServerError)
		return
	}
	responser.SendOK(w, "Сообщение удалено", http.StatusOK)
}

// UpdateMessage godoc
// @Summary Update message
// @Tags message
// @Param message body models.MessageInput true "Message info"
// @Param messageId path string true "messageId ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Сообщение успешно изменено"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} customerror.NoPermissionError "Нет доступа"
// @Failure 500	{object} responser.ErrorResponse "Не удалось обновить сообщение"
// @Router /messages/{messageId} [put]
func (h *MessageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "UpdateMessage")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	messageId := mapVars["messageId"]
	log.Printf("messageId: %s", messageId)
	messageUUID, err := uuid.Parse(messageId)
	if err != nil {
		//conn.400
		log.Printf("Получен кривой Id сообщения %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Получен кривой Id сообщения %v", err), http.StatusBadRequest)
		return
	}
	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		log.Println("нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	var messageDTO models.Message
	err = json.NewDecoder(r.Body).Decode(&messageDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Не удалось распарсить Json: %v", err), http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateMessage(ctx, user, messageUUID, messageDTO)

	if err != nil {
		if errors.As(err, &noPerm) {
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("внутренняя ошибка: %v", err), http.StatusInternalServerError)
		return
	}
	responser.SendOK(w, "Сообщение обновлено", http.StatusOK)
}


// GetAllMessages godoc
// @Summary Get All messages
// @Tags message
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param message body models.MessagesArrayDTO true "Messages"
// @Success 200 {object} models.MessagesArrayDTO "Сообщение успешно отаправлены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось получить сообщениея"
// @Router /chat/{chatId}/messages [get]
func (h *MessageController) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "GetAllMessages")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	chatId := mapVars["chatId"]
	log.Printf("chatid: %s", chatId)
	chatUUID, err := uuid.Parse(chatId)

	if err != nil {
		//conn.400
		log.Printf("Получен кривой Id юзера: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Получен кривой Id юзера: %v", err), http.StatusBadRequest)
		return
	}

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting getting all messages for chat: %v", chatUUID)

	messages, err := h.usecase.GetFirstMessages(r.Context(), chatUUID)
	if err != nil {
		log.Println("Error reading message:", err)
		responser.SendError(ctx, w, fmt.Sprintf("Error reading message:%v", err), http.StatusInternalServerError)
		return
	}

	if err := validator.Check(messages); err != nil {
		log.Printf("выходные данные не прошли проверку валидации: %v", err)
		responser.SendError(ctx, w, "Invalid data", http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(messages)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		responser.SendError(ctx, w, fmt.Sprintf("error happened in JSON marshal. Err: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}


// GetMessagesWithPage godoc
// @Summary получить 25 сообщений до определенного
// @Tags message
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param lastMessageId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} models.MessagesArrayDTO "Сообщение успешно отаправлены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} customerror.NoPermissionError "Нет доступа"
// @Failure 500	{object} responser.ErrorResponse "Не удалось получить сообщениея"
// @Router /chat/{chatId}/messages/pages/{lastMessageId} [get]
func (h *MessageController) GetMessagesWithPage(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "GetMessagesWithPage")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)

	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}
	chatId := mapVars["chatId"]
	lastMessageId := mapVars["lastMessageId"]
	log.Printf("chatid: %s, lastMessageId: %v", chatId, lastMessageId)

	chatUUID, err := uuid.Parse(chatId)
	if err != nil {
		//conn.400
		log.Printf("получен кривой Id юзера: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Delivery: error during parsing uuid:%v", err), http.StatusBadRequest)
		return
	}

	lastMessageUUID, err := uuid.Parse(lastMessageId)
	if err != nil {
		//conn.400
		log.Printf("получен кривой Id сообщения: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Delivery: error during parsing uuid:%v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	log.Println(user)
	if !ok {
		log.Println("нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	messages, err := h.usecase.GetMessagesWithPage(ctx, user.ID, chatUUID, lastMessageUUID)
	if err != nil {
		if errors.As(err, &noPerm) {
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("внутренняя ошибка: %v", err), http.StatusInternalServerError)
		return
	}

	responser.SendStruct(ctx, w, messages, http.StatusOK)
}

// SearchMessages godoc
// @Summary поиск сообщений
// @Tags message
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param search_query query int false "Поиск" example(sosal?)
// @Success 200 {object} models.MessagesArrayDTO "Сообщение успешно отаправлены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 403	{object} customerror.NoPermissionError "Нет доступа"
// @Failure 500	{object} responser.ErrorResponse "Не удалось получить сообщениея"
// @Router /chat/{chatId}/messages/search [get]
func (h *MessageController) SearchMessages(w http.ResponseWriter, r *http.Request) {
	metric.IncHit()
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestMessageDuration, "SearchMessages")
	}()

	log := logger.LoggerWithCtx(r.Context(), logger.Log)

	ctx := r.Context()
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}
	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)
	if err != nil {
		//conn.400
		log.Printf("получен кривой Id юзера: %v", err)
		responser.SendError(ctx, w, fmt.Sprintf("Delivery: error during parsing uuid:%v", err), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	log.Println(user)
	if !ok {
		log.Println("нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}
	log.Println(r.URL.Query())
	query := r.URL.Query().Get("search_query")
	if query == "" {
		log.Errorf("Поисковый запрос пуст")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	messages, err := h.usecase.SearchMessagesWithQuery(ctx, user, chatUUID, query)

	if err != nil {
		if errors.As(err, &noPerm) {
			responser.SendError(ctx, w, fmt.Sprintf("Нет доступа: %v", err), http.StatusForbidden)
			return
		}
		responser.SendError(ctx, w, fmt.Sprintf("внутренняя ошибка: %v", err), http.StatusInternalServerError)
		return
	}
	responser.SendStruct(ctx, w, messages, http.StatusOK)
}
