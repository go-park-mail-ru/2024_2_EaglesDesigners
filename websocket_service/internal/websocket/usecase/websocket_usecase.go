package usecase

import (
	"context"
	"net"
	"strconv"

	chatModels "github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/events"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	grpcChat "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ивент может быть либо изменение сущности чата, либо сообщение
type AnyEvent struct {
	TypeOfEvent string
	Event       interface{}
}

type ChatInfo struct {
	events chan chatModels.Event
	users  map[uuid.UUID]struct{}
}

type WebsocketUsecase struct {
	ch *amqp.Channel
	// мапа с чатами и каналами для ивентов по чатам
	onlineChats map[uuid.UUID]ChatInfo
	// мапа с онлайн пользователями и
	onlineUsers    map[uuid.UUID]chan AnyEvent
	chatRepository grpcChat.ChatServiceClient
}

func NewWebsocketUsecase(ch *amqp.Channel, host string, port int) *WebsocketUsecase {
	_, err := ch.QueueDeclare(
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

	_, err = ch.QueueDeclare(
		"chat", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log := logger.LoggerWithCtx(context.Background(), logger.Log)
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	grpcAddress := net.JoinHostPort(host, strconv.Itoa(port))
	// Создаем клиент
	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		// Используем insecure-коннект для тестов
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("s")
	}

	// gRPC-клиент сервера Auth
	grpcClient := grpcChat.NewChatServiceClient(cc)

	socket := &WebsocketUsecase{
		ch:             ch,
		onlineChats:    map[uuid.UUID]ChatInfo{},
		onlineUsers:    map[uuid.UUID]chan AnyEvent{},
		chatRepository: grpcClient,
	}

	go socket.consumeMessages()
	go socket.consumeChats()

	return socket
}

func (w *WebsocketUsecase) InitBrokersForUser(userId uuid.UUID, eventChannel chan AnyEvent) error {
	log := logger.LoggerWithCtx(context.Background(), logger.Log)

	chats, err := w.chatRepository.GetUserChats(context.Background(), &grpcChat.UserChatsRequest{UserId: userId.String()})
	if err != nil {
		log.Errorf("Не удалось запросить чаты: %v", err)
		return err
	}

	w.onlineUsers[userId] = eventChannel
	log.Infof("Пользователь %v онлайн", userId)

	// Добавляем в брокеры пользователей
	for _, chatId := range chats.ChatIds {
		chatUUID, err := uuid.Parse(chatId)
		if err != nil {
			continue
		}

		if chatInfo, ok := w.onlineChats[chatUUID]; ok {
			chatInfo.events <- chatModels.Event{
				Action: AddWebcosketUser,
				Users:  []uuid.UUID{userId},
			}
			log.Infof("Пользователь %v добавлен в брокер для чата %v", userId, chatId)
		}
	}
	return nil
}
