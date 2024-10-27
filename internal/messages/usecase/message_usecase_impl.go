package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"

	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

type MessageUsecaseImplm struct {
	messageRepository repository.MessageRepository
	tokenUsecase      *usecase.Usecase
	redisClient       *redis.Client
	messages          chan models.Message
	activeUsers       map[uuid.UUID]bool
}

func NewMessageUsecaseImpl(messageRepository repository.MessageRepository, tokenUsecase *usecase.Usecase, redisClient *redis.Client) MessageUsecase {
	usecase := MessageUsecaseImplm{
		messageRepository: messageRepository,
		tokenUsecase:      tokenUsecase,
		redisClient:       redisClient,
		messages:          make(chan models.Message, 100),
		activeUsers:       map[uuid.UUID]bool{},
	}
	go usecase.goBroker(context.Background())
	return &usecase
}

func (u *MessageUsecaseImplm) publishMessageIvent(ctx context.Context, message models.Message) error {
	log.Println("Message usecase: добавление сообщения в redis")

	err := u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: message.AuthorID.String(),
		MaxLen: 0,
		ID:     "",
		Values: message,
	}).Err()

	return err
}

func (u *MessageUsecaseImplm) goBroker(ctx context.Context) {
	for {
		select {
		case message := <-u.messages:
			if ok := u.activeUsers[message.AuthorID]; ok {
				err := u.publishMessageIvent(ctx, message)
				if err != nil {
					log.Printf("Message usecase: не удалось отправить в поток: %v", err)
				}

			}
		default:
		}
	}
}

func (u *MessageUsecaseImplm) SendMessage(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, message models.Message) error {
	log.Printf("Usecase: начато добавление сообщения в чат %v", chatId)

	message.MessageId = uuid.New()
	message.SentAt = time.Now()
	message.AuthorID = userId

	log.Printf("Usecase: сообщение от прользователя: %v", message.AuthorID)

	err := u.messageRepository.AddMessage(message, chatId)
	if err != nil {
		log.Printf("Usecase: не удалось добавить сообщение: %v", err)
		return err
	}

	// записываем новое сообщение в канал
	u.messages <- message

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

func (u *MessageUsecaseImplm) ScanForNewMessages(ctx context.Context, channel chan<- []models.Message, chatId uuid.UUID, res chan<- error, closeChannel <-chan bool) {
	defer func() {
		close(channel)
		close(res)
	}()
	user, ok := ctx.Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		return
	}

	u.activeUsers[user.ID] = true
	defer func() { u.activeUsers[user.ID] = false }()

	log.Println("Message usecase: начат поиск новых сообщений")

	duration := 500 * time.Millisecond

	for {
		select {
		case <-closeChannel:
			log.Println("Message usecase: scanning stoped")
			return
		default:
			time.Sleep(duration)
			// Чтение сообщений из Stream
			messages, err := u.redisClient.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{user.ID.String(), "0"}, // Начинаем с самого начала ("0")
				Count:   5,                               // Получить 5 сообщений
				Block:   0,                               // Блокировать до появления новых сообщений
			}).Result()

			// получаем новые сообщения в канал
			for _, message := range messages {
				fmt.Println("Стрим:", message.Stream)
				for _, msg := range message.Messages {
					fmt.Printf("ID: %s, Данные: %v\n", msg.ID, msg.Values)
				}
			}

			if err != nil {
				fmt.Println("Ошибка при чтении сообщений:", err)
			}

		}
	}
}
