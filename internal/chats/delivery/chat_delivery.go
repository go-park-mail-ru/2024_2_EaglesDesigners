package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/usecase"
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
// @Success 200 "Пользователи добавлены"
// @Failure 400	"Некорректный запрос"
// @Failure 500	"Не удалось добавить пользователей"
// @Router /addusers [post]
func (c *ChatDelivery) AddUsersIntoChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usersToAdd model.AddUsersIntoChatDTO
	err := json.NewDecoder(r.Body).Decode(&usersToAdd)
	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.AddUsersIntoChat(context.Background(), r.Cookies(), usersToAdd.UsersId, usersToAdd.ChatId)

	if err != nil {
		log.Printf("Не удалось добавить пользователей в чат: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}
