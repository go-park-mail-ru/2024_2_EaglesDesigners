package usecase

import (
	"context"
	"fmt"
	"time"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"

	"github.com/google/uuid"
)

type Method = string

const (
	FeatNewUser Method = "featNewUser"
	DelUser     Method = "deleteUser"
	NewMessage  Method = "message"
)

type MessageUsecaseImplm struct {
	messageRepository repository.MessageRepository
	chatRepository    chatRepository.ChatRepository
	tokenUsecase      *usecase.Usecase
}

func NewMessageUsecaseImpl(messageRepository repository.MessageRepository, chatRepository chatRepository.ChatRepository, tokenUsecase *usecase.Usecase) MessageUsecase {
	usecase := MessageUsecaseImplm{
		messageRepository: messageRepository,
		tokenUsecase:      tokenUsecase,
		chatRepository:    chatRepository,
	}
	return &usecase
}

func (u *MessageUsecaseImplm) SendMessage(ctx context.Context, user jwt.User, chatId uuid.UUID, message models.Message) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Usecase: начато добавление сообщения в чат %v", chatId)

	message.MessageId = uuid.New()
	message.SentAt = time.Now()
	message.AuthorID = user.ID
	message.ChatId = chatId

	log.Printf("Usecase: сообщение от прользователя: %v", message.AuthorID)

	err := u.messageRepository.AddMessage(message, chatId)
	if err != nil {
		log.Errorf("Usecase: не удалось добавить сообщение: %v", err)
		return err
	}

	log.Printf("Usecase: сообщение успешно добавлено: %v", message.MessageId)
	return nil
}

func (u *MessageUsecaseImplm) GetFirstMessages(ctx context.Context, chatId uuid.UUID) (models.MessagesArrayDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Usecase: начато получение сообщений")

	messages, err := u.messageRepository.GetFirstMessages(ctx, chatId)
	if err != nil {
		log.Errorf("Usecase: не удалось получить сообщения: %v", err)
		return models.MessagesArrayDTO{}, err
	}
	log.Printf("Usecase: сообщения получены")

	return models.MessagesArrayDTO{
		Messages: messages,
	}, nil
}

func (u *MessageUsecaseImplm) GetMessagesWithPage(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, lastMessageId uuid.UUID) (models.MessagesArrayDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("запрошены сообщения из чата: %v, запрос получен от пользовтеля: %v", chatId, userId)

	_, err := u.chatRepository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		log.Printf("пользователь %v не состоит в чате %v", userId, chatId)
		return models.MessagesArrayDTO{},
			&customerror.NoPermissionError{
				Area: fmt.Sprintf("чат %v", chatId),
				User: fmt.Sprintf("пользователь %v", userId),
			}
	}

	messages, err := u.messageRepository.GetAllMessagesAfter(ctx, chatId, lastMessageId)
	if err != nil {
		return models.MessagesArrayDTO{}, err
	}

	return models.MessagesArrayDTO{
		Messages: messages,
	}, nil

}
