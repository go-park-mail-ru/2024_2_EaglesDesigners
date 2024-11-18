package usecase

import (
	"context"
	"fmt"
	"time"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	socketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/websocket/usecase"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
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
	tokenUsecase      *jwt.Usecase
	queryName         string
	ch                *amqp.Channel
}

func NewMessageUsecaseImpl(messageRepository repository.MessageRepository, chatRepository chatRepository.ChatRepository, tokenUsecase *jwt.Usecase, ch *amqp.Channel) MessageUsecase {
	// объявляем очередь
	q, err := ch.QueueDeclare(
		"message", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log := logger.LoggerWithCtx(context.Background(), logger.Log)
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	usecase := MessageUsecaseImplm{
		messageRepository: messageRepository,
		tokenUsecase:      tokenUsecase,
		chatRepository:    chatRepository,
		queryName:         q.Name,
		ch:                ch,
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
	u.sendIvent(ctx, socketUsecase.NewMessage, message)
	return nil
}

func (u *MessageUsecaseImplm) DeleteMessage(ctx context.Context, user jwt.User, messageId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Infof("начато удаление сообщения %v пользователем %v", messageId, user.ID)

	message, err := u.messageRepository.GetMessageById(ctx, messageId)
	if err != nil {
		return err
	}

	if user.ID != message.AuthorID {
		return &customerror.NoPermissionError{
			Area: fmt.Sprintf("сообщение %v принадлежит другому пользователю", messageId),
			User: user.ID.String(),
		}
	}
	err = u.messageRepository.DeleteMessage(ctx, messageId)
	if err != nil {
		return err
	}

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

func (s *MessageUsecaseImplm) sendIvent(ctx context.Context, action string, message models.Message) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	newEvent := socketUsecase.MessageEvent{
		Action:  action,
		Message: message,
	}

	body, err := socketUsecase.SerializeMessageEvent(newEvent)
	if err != nil {
		log.Errorf("Не удалось сериализовать объект")
		return
	}
	err = s.ch.PublishWithContext(ctx,
		"",          // exchange
		s.queryName, // имя очереди
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("failed to publish a message. Error: %s", err)
	}
}
