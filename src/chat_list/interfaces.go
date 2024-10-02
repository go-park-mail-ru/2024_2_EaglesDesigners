package chatlist

import (
	"net/http"

	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
)

type ChatRepository interface {
	GetUserChats(user *userModel.User) []chatModel.Chat
}

type ChatService interface {
	GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error)
}
