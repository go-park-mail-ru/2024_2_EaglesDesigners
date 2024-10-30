package repository

import (
	"context"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/google/uuid"
)

type ChatRepository interface {
	GetUserChats(ctx context.Context, userId uuid.UUID, pageNum int) (chats []chatModel.Chat, err error)
	GetUserRoleInChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) (string, error)
	CreateNewChat(ctx context.Context, chat chatModel.Chat) error
	AddUserIntoChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, userROle string) error
	GetCountOfUsersInChat(ctx context.Context, chatId uuid.UUID) (int, error)
	GetChatById(ctx context.Context, chatId uuid.UUID) (chatModel.Chat, error)
}
