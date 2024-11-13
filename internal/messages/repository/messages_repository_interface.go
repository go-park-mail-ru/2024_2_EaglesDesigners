package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/google/uuid"
)

type MessageRepository interface {
	GetMessages(chatId uuid.UUID) ([]models.Message, error)
	AddMessage(message models.Message, chatId uuid.UUID) error
	GetLastMessage(chatId uuid.UUID) (models.Message, error)
	GetAllMessagesAfter(ctx context.Context, chatId uuid.UUID, lastMessageId uuid.UUID) ([]models.Message, error)
}
