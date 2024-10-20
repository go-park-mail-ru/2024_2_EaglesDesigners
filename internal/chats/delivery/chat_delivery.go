package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
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
func (c *ChatDelivery) GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Пришёл запрос на получения чатов с параметрами: %v", r.URL.Query())
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		log.Printf("Неверно указан параметр запроса page. page = %s. ERROR: %v", r.URL.Query().Get("page"), err)
	}

	chats, err := c.service.GetChats(context.Background(), r.Cookies(), pageNum)

	if err != nil {
		fmt.Println(err)

		//вернуть 401
		w.WriteHeader(http.StatusUnauthorized)

		log.Printf("НЕ УДАЛОСЬ ПОЛУЧИТЬ ЧАТЫ. ОШИБКА: %s", err)
		return
	}

	chatsDTO := models.ChatsDTO{
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

func (c *ChatDelivery) AddNewChat(w http.ResponseWriter, r *http.Request) {
	var chatDTO models.ChatDTO
	err := json.NewDecoder(r.Body).Decode(&chatDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.AddNewChat(context.Background(), r.Cookies(), chatDTO)
	if err != nil {
		log.Printf("Не удалось добавить чат: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (c *ChatDelivery) AddUsersIntoChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usersToAdd models.AddUsersIntoChatDTO
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
	return
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}
