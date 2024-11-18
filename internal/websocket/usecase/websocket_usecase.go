package usecase

import (
	"context"

	chatEvent "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatModels "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ивент может быть либо изменение сущности чата, либо сообщение
type AnyEvent struct {
	TypeOfEvent string
	Event       interface{}
}

type ChatInfo struct {
	events chan chatEvent.Event
	users  map[uuid.UUID]struct{}
}

type WebsocketUsecase struct {
	ch *amqp.Channel
	// мапа с чатами и каналами для ивентов по чатам
	onlineChats map[uuid.UUID]ChatInfo
	// мапа с онлайн пользователями и
	onlineUsers    map[uuid.UUID]chan AnyEvent
	chatRepository chatRepository.ChatRepository
}

func NewWebsocketUsecase(ch *amqp.Channel) *WebsocketUsecase {

	socket := &WebsocketUsecase{
		ch:          ch,
		onlineChats: map[uuid.UUID]ChatInfo{},
		onlineUsers: map[uuid.UUID]chan AnyEvent{},
	}

	go socket.consumeMessages()
	go socket.consumeChats()

	return socket
}

func (w *WebsocketUsecase) InitBrokersForUser(userId uuid.UUID, eventChannel <-chan AnyEvent) error {
	log := logger.LoggerWithCtx(context.Background(), logger.Log)

	chats, err := w.chatRepository.GetUserChats(context.Background(), userId)
	if err != nil {
		return err
	}

	w.onlineUsers[userId] = make(chan AnyEvent, 10)

	// Добавляем в брокеры пользователей
	for _, chat := range chats {
		if chatInfo, ok := w.onlineChats[chat.ChatId]; ok {
			chatInfo.events <- chatModels.Event{
				Action: AddWebcosketUser,
				Users:  []uuid.UUID{userId},
			}
			log.Debugf("Пользователь %v добавлен в броке для чата %v", userId, chat.ChatId)
		}
	}
	return nil
}
