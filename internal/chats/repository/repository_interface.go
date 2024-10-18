package repository

import (
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/google/uuid"
)

type ChatRepository interface {
	GetUserChats(userId uuid.UUID, pageSize int) (chats []chatModel.ChatDAO, err error)
	IsUserInChat(userId uuid.UUID, chatId uuid.UUID) (bool, error)
	CreateNewChat(chat chatModel.Chat) error
	AddUserIntoChat(userId uuid.UUID, chatId uuid.UUID, userROle string) error
}
