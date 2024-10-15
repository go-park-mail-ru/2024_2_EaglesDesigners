package repository

import (
	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
)

type ChatRepository interface {
	GetUserChats(user *userModel.User) []chatModel.Chat
}
