package usecase

import (
	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
)

type ChatUsecase interface {
	GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error)
}
