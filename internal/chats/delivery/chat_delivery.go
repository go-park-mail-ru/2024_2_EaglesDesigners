package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
func (c *ChatDelivery) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Пришёл запрос на получения чатов")

	_, err := c.service.GetChats(context.Background(), r.Cookies(), 0)
	if err != nil {
		fmt.Println(err)

		//вернуть 401
		w.WriteHeader(http.StatusUnauthorized)

		log.Printf("НЕ УДАЛОСЬ ПОЛУЧИТЬ ЧАТЫ. ОШИБКА: %s", err)
		return
	}

	chatsDTO := models.ChatsDTO{
		Chats: nil,
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
