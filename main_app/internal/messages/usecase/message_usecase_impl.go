package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	socketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/events"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	"github.com/prometheus/client_golang/prometheus"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/custom_error"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/repository"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Method = string

func init() {
	prometheus.MustRegister(sendedMessagesMetric, deleteMessageMetric, updateMessageMetric)
}

const (
	FeatNewUser Method = "featNewUser"
	DelUser     Method = "deleteUser"
	NewMessage  Method = "message"
)

type MessageUsecaseImplm struct {
	fileUC            FilesUsecase
	messageRepository repository.MessageRepository
	chatRepository    chatRepository.ChatRepository
	queryName         string
	ch                *amqp.Channel
}

type FilesUsecase interface {
	SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (string, error)
}

func NewMessageUsecaseImpl(fileUC FilesUsecase, messageRepository repository.MessageRepository, chatRepository chatRepository.ChatRepository, ch *amqp.Channel) MessageUsecase {
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
		fileUC:            fileUC,
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
		queryName:         q.Name,
		ch:                ch,
	}
	return &usecase
}

var sendedMessagesMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_sended_messages",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

const branchType = "branch"

func (u *MessageUsecaseImplm) SendMessage(ctx context.Context, user jwt.User, chatId uuid.UUID, message models.Message) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Usecase: начато добавление сообщения в чат %v", chatId)

	message.MessageId = uuid.New()
	message.SentAt = time.Now()
	message.AuthorID = user.ID
	message.ChatId = chatId

	log.Printf("Usecase: сообщение от прользователя: %v", message.AuthorID)

	chatType, err := u.chatRepository.GetChatType(ctx, chatId)
	if err != nil {
		log.Errorf("Usecase: не удалось получить тип чата: %v", err)
		return err
	}

	// Добавление тех, кому можно читать файл, если файл не в публичное место отправлен.
	var userIDs []string
	if chatType == "personal" || chatType == "group" {
		users, err := u.chatRepository.GetUsersFromChat(ctx, chatId)
		if err != nil {
			log.Errorf("Usecase: не удалось получить пользователей чата: %v", err)
			return err
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID.String())
		}
	}

	// Добавление файлов в mongoDB
	for i := 0; i < len(message.Files); i++ {
		fileURl, err := u.fileUC.SaveFile(ctx, message.Files[i], message.FilesHeaders[i], userIDs)
		if err != nil {
			log.Errorf("Usecase: не удалось сохранить файл: %v", err)
			return err
		}
		message.FilesURLs = append(message.FilesURLs, fileURl)
	}

	err = u.messageRepository.AddMessage(message, chatId)
	if err != nil {
		log.Errorf("Usecase: не удалось добавить сообщение: %v", err)
		return err
	}

	log.Printf("Usecase: сообщение успешно добавлено: %v", message.MessageId)

	// Если это ветка, то нужно указать parent branch.
	chatType, err := u.chatRepository.GetChatType(ctx, chatId)
	if err != nil {
		chatType = ""
	}
	if chatType == branchType {
		parentChatId, err := u.chatRepository.GetBranchParent(ctx, chatId)
		if err != nil {
			return err
		}
		message.ChatIdParent = parentChatId
	}
	u.SendIvent(ctx, socketUsecase.NewMessage, message)

	metric.IncMetric(*sendedMessagesMetric)
	return nil
}

var deleteMessageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_deleted_messages",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

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

	u.SendIvent(ctx, socketUsecase.DeleteMessage, message)
	metric.IncMetric(*deleteMessageMetric)
	return nil
}

var updateMessageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_updated_messages",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

// не нужен, можно в repo.AddMessage тип определять по наличию файлов
func (u *MessageUsecaseImplm) SendInformationalMessage(_ context.Context, message models.Message, chatId uuid.UUID) error {
	return u.messageRepository.AddInformationalMessage(message, chatId)
}

func (u *MessageUsecaseImplm) UpdateMessage(ctx context.Context, user jwt.User, messageId uuid.UUID, message models.Message) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Infof("Начато изменение сообщения %v. Запрос от пользователя %v", messageId, user.ID)

	newText := message.Message

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

	u.messageRepository.UpdateMessage(ctx, messageId, newText)

	// отправляем в сокет
	message.Message = newText
	message.IsRedacted = true
	u.SendIvent(ctx, socketUsecase.UpdateMessage, message)
	metric.IncMetric(*updateMessageMetric)
	return nil
}

const NotInChat = ""

func (u *MessageUsecaseImplm) SearchMessagesWithQuery(ctx context.Context, user jwt.User, chatId uuid.UUID, searchQuery string) (models.MessagesArrayDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Infof("Начат поиск сообщений в чате %v. Поисковая строка = %v", chatId, searchQuery)

	role, err := u.chatRepository.GetUserRoleInChat(ctx, user.ID, chatId)
	if err != nil {
		return models.MessagesArrayDTO{}, err
	}

	if role == NotInChat {
		return models.MessagesArrayDTO{},
			&customerror.NoPermissionError{
				Area: fmt.Sprintf("Нет доступа к чату %v", chatId),
				User: user.ID.String(),
			}
	}

	messages, err := u.messageRepository.SearchMessagesWithQuery(ctx, chatId, searchQuery)

	if err != nil {
		return models.MessagesArrayDTO{}, err
	}

	return models.MessagesArrayDTO{
		Messages: messages,
	}, nil
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

	role, err := u.chatRepository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return models.MessagesArrayDTO{}, err
	}

	if role == NotInChat {
		log.Printf("пользователь %v не состоит в чате %v", userId, chatId)
		return models.MessagesArrayDTO{},
			&customerror.NoPermissionError{
				Area: fmt.Sprintf("чат %v", chatId),
				User: fmt.Sprintf("пользователь %v", userId),
			}
	}

	messages, err := u.messageRepository.GetMessagesAfter(ctx, chatId, lastMessageId)
	if err != nil {
		return models.MessagesArrayDTO{}, err
	}

	return models.MessagesArrayDTO{
		Messages: messages,
	}, nil

}

func (s *MessageUsecaseImplm) GetLastMessage(chatId uuid.UUID) (models.Message, error) {
	return s.messageRepository.GetLastMessage(chatId)
}

func (s *MessageUsecaseImplm) SendIvent(ctx context.Context, action string, message models.Message) {
	newMessage := socketUsecase.Message{
		MessageId:  message.MessageId,
		AuthorID:   message.AuthorID,
		BranchID:   message.BranchID,
		Message:    message.Message,
		SentAt:     message.SentAt,
		ChatId:     message.ChatId,
		IsRedacted: message.IsRedacted,
	}

	log := logger.LoggerWithCtx(ctx, logger.Log)
	newEvent := socketUsecase.MessageEvent{
		Action:  action,
		Message: newMessage,
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
