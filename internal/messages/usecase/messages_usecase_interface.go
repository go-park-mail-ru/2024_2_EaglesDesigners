package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/google/uuid"
)

type MessageUsecase interface {
	SendMessage(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, message models.Message) error
	GetMessages(ctx context.Context, chatId uuid.UUID, pageId int) (models.MessagesArrayDTO, error)
	ScanForNewMessages(ctx context.Context, channel chan<- []models.Message, chatId uuid.UUID, res chan<- error, closeChannel <-chan bool)
}
