package usecase

import (
	"context"

	chatEvent "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
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

type WebsocketUsecase struct {
	ch             *amqp.Channel
	// мапа с чатами и каналами для ивентов по чатам
	onlineChats    map[uuid.UUID]chan chatEvent.Event
	// мапа с онлайн пользователями и 
	onlineUsers    map[uuid.UUID]chan AnyEvent
	chatRepository chatRepository.ChatRepository
}

func NewWebsocketUsecase(ch *amqp.Channel) *WebsocketUsecase {

	socket := &WebsocketUsecase{
		ch:          ch,
		onlineChats: map[uuid.UUID]chan chatEvent.Event{},
		onlineUsers: map[uuid.UUID]chan AnyEvent{},
	}

	go socket.consumeMessages()
	go socket.consumeChats()

	return socket
}



func (w *WebsocketUsecase) InitBrokersForUser(userId uuid.UUID, eventChannel <- chan AnyEvent) error {
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	
	chats, err := w.chatRepository.GetUserChats(context.Background(), userId, 0)
}