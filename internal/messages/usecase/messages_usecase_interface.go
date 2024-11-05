package usecase

import (
	"context"

	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"

	"github.com/google/uuid"
)

type MessageUsecase interface {
	SendMessage(ctx context.Context, user jwt.User, chatId uuid.UUID, message models.Message) error
	GetMessages(ctx context.Context, chatId uuid.UUID) (models.MessagesArrayDTO, error)
	ScanForNewMessages(ctx context.Context, channel chan<- models.WebScoketDTO, res chan<- error, closeChannel <-chan bool)
	GetOnlineUsers() map[uuid.UUID]bool
}
