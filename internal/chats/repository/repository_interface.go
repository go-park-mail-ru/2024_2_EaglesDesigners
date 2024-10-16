package repository

import (
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type ChatRepository interface {
	GetUserChats(user *userModel.User) []chatModel.Chat
	IsUserInChat(userId int, chatId int) bool
}
