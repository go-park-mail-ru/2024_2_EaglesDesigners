package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/usecase"
)

type ChatController struct {
	service chatlist.ChatUsecase
}

func NewChatController(service chatlist.ChatUsecase) *ChatController {
	return &ChatController{
		service: service,
	}
}

// ChatHandler godoc
// @Summary Get user chats
// @Description Retrieve the list of chats for the authenticated user based on their access token.
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ChatsDTO "List of chats"
// @Failure 401 {object} ErrorResponse "Unauthorized, no valid access token"
// @Router /chats [get]
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

type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}
