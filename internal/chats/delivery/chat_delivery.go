package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/usecase"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ChatDelivery struct {
	service chatlist.ChatUsecase
}

func NewChatDelivery(service chatlist.ChatUsecase) *ChatDelivery {
	return &ChatDelivery{
		service: service,
	}
}

// GetUserChatsHandler выдает чаты пользователя в query указать страницу ?page=
//
// GetUserChatsHandler godoc
// @Summary Get chats of user
// @Tags chat
// @Produce json
// @Param page query int false "Page number for pagination" default(0)
// @Success 200 {object} model.ChatsDTO
// @Failure 500	"Не удалось получить сообщения"
// @Router /chats [get]
func (c *ChatDelivery) GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Пришёл запрос на получения чатов с параметрами: %v", r.URL.Query())
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		log.Printf("Неверно указан параметр запроса page. page = %s. ERROR: %v", r.URL.Query().Get("page"), err)
		pageNum = 0
	}

	chats, err := c.service.GetChats(r.Context(), r.Cookies(), pageNum)

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

	jsonResp, err := json.Marshal(chatsDTO)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// AddNewChat godoc
// @Summary Add new chat
// @Tags chat
// @Accept json
// @Param chat body model.ChatDTOInput true "Chat info"
// @Success 201 "Чат создан"
// @Failure 400	"Некорректный запрос"
// @Failure 500	"Не удалось добавить чат / группу"
// @Router /addchat [post]
func (c *ChatDelivery) AddNewChat(w http.ResponseWriter, r *http.Request) {
	var chatDTO model.ChatDTOInput
	err := json.NewDecoder(r.Body).Decode(&chatDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.AddNewChat(r.Context(), r.Cookies(), chatDTO)
	if err != nil {
		log.Printf("Не удалось добавить чат: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddUsersIntoChat godoc
// @Summary Добавить пользователей в чат
// @Tags chat
// @Accept json
// @Param users body model.AddUsersIntoChatDTO true "Пользователи на добавление"
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Пользователи добавлены"
// @Failure 400	"Некорректный запрос"
// @Failure 500	"Не удалось добавить пользователей"
// @Router /chat/{chatId}/addusers [post]
func (c *ChatDelivery) AddUsersIntoChat(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> AddUsersIntoChat: error parsing chat uuid:", err)
	}

	var usersToAdd model.AddUsersIntoChatDTO
	err = json.NewDecoder(r.Body).Decode(&usersToAdd)
	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.AddUsersIntoChat(r.Context(), r.Cookies(), usersToAdd.UsersId, chatUUID)

	if err != nil {
		log.Printf("Не удалось добавить пользователей в чат: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteChatOrGroup godoc
// @Summary Удаличть чат или группу
// @Tags chat
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 "Чат удалён"
// @Failure 400	"Некорректный запрос"
// @Failure 403	"Нет полномочий"
// @Failure 500	"Не удалось удалить чат"
// @Router /chat/{chatId}/delete [delete]
func (c *ChatDelivery) DeleteChatOrGroup(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := getChatIdFromContext(r.Context())

	if err != nil {
		//conn.400
		log.Println("Chat delivery -> AddUsersIntoChat: error parsing chat uuid:", err)
	}

	user, ok := r.Context().Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.service.DeleteChat(r.Context(), chatUUID, user.ID)

	if err != nil {
		if errors.As(err, &customerror.NoPermissionError{}) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getChatIdFromContext(ctx context.Context) (uuid.UUID, error) {
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
