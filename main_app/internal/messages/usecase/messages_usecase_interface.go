package usecase

import (
	"context"

	"github.com/google/uuid"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
)

//go:generate mockgen -source=messages_usecase_interface.go -destination=mocks/mocks.go

type MessageUsecase interface {
	SendMessage(ctx context.Context, user auth.User, chatId uuid.UUID, message models.Message) error
	SendInformationalMessage(ctx context.Context, message models.Message, chatId uuid.UUID) error

	DeleteMessage(ctx context.Context, user auth.User, messageId uuid.UUID) error
	UpdateMessage(ctx context.Context, user auth.User, messageId uuid.UUID, message models.Message) error

	SearchMessagesWithQuery(ctx context.Context, user auth.User, chatId uuid.UUID, searchQuery string) (models.MessagesArrayDTO, error)
	GetMessagesWithPage(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, lastMessageId uuid.UUID) (models.MessagesArrayDTO, error)

	GetFirstMessages(ctx context.Context, chatId uuid.UUID) (models.MessagesArrayDTO, error)
	GetLastMessage(chatId uuid.UUID) (models.Message, error)

	SendIvent(ctx context.Context, action string, message models.Message)
	GetPayload(ctx context.Context, chatId uuid.UUID) (files []models.Payload, photos []models.Payload, err error)
}
