package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/service"
)

type ChatController struct {
	service service.ChatService
}

func NewChatController(service service.ChatService) *ChatController {
	return &ChatController{
		service: service,
	}
}

func (c *ChatController) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Пришёл запрос на получения чатов")

	chats, err := c.service.GetChats(r.Cookies())
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
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
