package usecase

import (
	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
)

type ChatUsecase interface {
	GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error)

	// CanUserWriteInChat проверяет может ли юзер писать в чат
	CanUserWriteInChat(userId int, chatId int) bool
}
