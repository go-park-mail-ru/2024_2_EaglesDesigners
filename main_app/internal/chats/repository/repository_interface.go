package repository

import (
	"context"

	"github.com/google/uuid"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
)

//go:generate mockgen -source=repository_interface.go -destination=mocks/mocks.go

type ChatRepository interface {
	// GetUserChats возвращает список чатов, в которых есть пользователь.
	GetUserChats(ctx context.Context, userId uuid.UUID) (chats []chatModel.Chat, err error)

	// если не состоит в чате, то вернет пустую строку
	GetUserRoleInChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) (string, error)
	CreateNewChat(ctx context.Context, chat chatModel.Chat) error
	AddUserIntoChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, userROle string) error
	GetCountOfUsersInChat(ctx context.Context, chatId uuid.UUID) (int, error)
	GetChatById(ctx context.Context, chatId uuid.UUID) (chatModel.Chat, error)
	GetChatType(ctx context.Context, chatId uuid.UUID) (string, error)
	DeleteChat(ctx context.Context, chatId uuid.UUID) error
	UpdateChat(ctx context.Context, chatId uuid.UUID, chatUpdate string) error
	DeleteUserFromChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error
	GetUsersFromChat(ctx context.Context, chatId uuid.UUID) ([]chatModel.UserInChatDAO, error)
	UpdateChatPhoto(ctx context.Context, chatId uuid.UUID, filename string) error
	GetNameAndAvatar(ctx context.Context, userId uuid.UUID) (string, string, error)
	SearchUserChats(ctx context.Context, userId uuid.UUID, keyWord string) ([]chatModel.Chat, error)
	SearchGlobalChats(ctx context.Context, userId uuid.UUID, keyWord string) ([]chatModel.Chat, error)

	// SetChatNotofications позволяет включить или выключить уведомления.
	SetChatNotofications(ctx context.Context, chatUUID uuid.UUID, userId uuid.UUID, value bool) error

	GetSendNotificationsForUser(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (bool, error)

	// GetBranchParent используется для веток, чтобы находить родителя.
	GetBranchParent(ctx context.Context, branchId uuid.UUID) (uuid.UUID, error)
	AddBranch(ctx context.Context, chatId uuid.UUID, messageId uuid.UUID) (chatModel.AddBranch, error)
}
