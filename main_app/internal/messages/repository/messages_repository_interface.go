package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=messages_repository_interface.go -destination=mocks/mocks.go

type MessageRepository interface {
	AddMessage(message models.Message, chatId uuid.UUID) error

	DeleteMessage(ctx context.Context, messageId uuid.UUID) error

	UpdateMessage(ctx context.Context, messageId uuid.UUID, newText string) error

	SearchMessagesWithQuery(ctx context.Context, chatId uuid.UUID, searchQuery string) ([]models.Message, error)
	GetFirstMessages(ctx context.Context, chatId uuid.UUID) ([]models.Message, error)
	GetMessageById(ctx context.Context, messageId uuid.UUID) (models.Message, error)
	GetLastMessage(chatId uuid.UUID) (models.Message, error)
	GetAllMessagesAfter(ctx context.Context, chatId uuid.UUID, lastMessageId uuid.UUID) ([]models.Message, error)
}
