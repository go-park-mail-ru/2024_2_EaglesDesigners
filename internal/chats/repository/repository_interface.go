package repository

import (
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/google/uuid"
)

type ChatRepository interface {
	GetUserChats(userId uuid.UUID, pageNum int) (chats []chatModel.Chat, err error)
	GetUserRoleInChat(userId uuid.UUID, chatId uuid.UUID) (string, error)
	CreateNewChat(chat chatModel.Chat) error
	AddUserIntoChat(userId uuid.UUID, chatId uuid.UUID, userROle string) error
}
