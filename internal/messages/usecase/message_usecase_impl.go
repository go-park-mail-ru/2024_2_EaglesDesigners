package usecase

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type MessageUsecaseImplm struct {
	messageRepository repository.MessageRepository
	tokenUsecase      *usecase.Usecase
}

func NewMessageUsecaseImpl(messageRepository repository.MessageRepository) MessageUsecase {
	return &MessageUsecaseImplm{
		messageRepository: messageRepository,
	}
}

func (u *MessageUsecaseImplm) SendMessage(ctx context.Context,  cookie []*http.Cookie, chatId uuid.UUID, message models.Message) error {
	log.Printf("Usecase: начато добавление сообщения в чат %v", chatId)

	user, err := u.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	log.Printf("Chat usecase: пришел запрос на получение всех чатов от пользователя: %v", user.ID)

	message.MessageId = uuid.New()
	message.SentAt = time.Now()
	message.AuthorID = user.ID

	log.Printf("Usecase: сообщение от прользователя: %v", message.AuthorID)

	err = u.messageRepository.AddMessage(message, chatId)
	if err != nil {
		log.Printf("Usecase: не удалось добавить сообщение: %v", err)
		return err
	}
	log.Printf("Usecase: сообщение успешно добавлено: %v", message.MessageId)

	return nil
}

func (u *MessageUsecaseImplm) GetMessages(ctx context.Context, chatId uuid.UUID, pageId int) (models.MessagesArrayDTO, error) {
	log.Printf("Usecase: начато получение сообщений")

	messages, err := u.messageRepository.GetMessages(pageId, chatId)
	if err != nil {
		log.Printf("Usecase: не удалось получить сообщения: %v", err)
		return models.MessagesArrayDTO{}, err
	}
	log.Printf("Usecase: сообщения получены")

	return models.MessagesArrayDTO{
		Messages: messages,
	}, nil
}

func (u *MessageUsecaseImplm) ScanForNewMessages(channel chan<- []models.Message, chatId uuid.UUID, res chan<- error, closeChannel <-chan bool) {
	defer func() {
		close(channel)
		close(res)
	}()

	startMessage, err := u.messageRepository.GetLastMessage(chatId)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		res <- err
	}

	duration := 500 * time.Millisecond

	for {

		if <-closeChannel {
			return
		}

		time.Sleep(duration)

		newMessage, err := u.messageRepository.GetLastMessage(chatId)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			res <- err
			continue
		}

		if newMessage.MessageId != startMessage.MessageId {
			messages, err := u.messageRepository.GetAllMessagesAfter(chatId, startMessage.SentAt, startMessage.MessageId)

			if err != nil {
				res <- err
				return
			}
			channel <- messages

			startMessage = newMessage
		}
	}
}
