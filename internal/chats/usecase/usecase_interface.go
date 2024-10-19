package usecase

import (
	"context"
	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/google/uuid"
)

type ChatUsecase interface {
	GetChats(ctx context.Context, cookie []*http.Cookie, pageNum int) ([]chatModel.Chat, error)
	AddUsersIntoChat(ctx context.Context, cookie []*http.Cookie, user_ids []uuid.UUID, chat_id uuid.UUID) error

	// CanUserWriteInChat проверяет может ли юзер писать в чат
	AddNewChat(ctx context.Context, cookie []*http.Cookie, chat chatModel.Chat) error
}
