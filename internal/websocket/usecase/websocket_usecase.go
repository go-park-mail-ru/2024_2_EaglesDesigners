package usecase

import (
	"log"

	chatEvent "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type WebsocketUsecase struct {
	ch          *amqp.Channel
	activeChats map[uuid.UUID]chan chatEvent.Event
}

func NewWebsocketUsecase(ch *amqp.Channel) *WebsocketUsecase {

	socket := &WebsocketUsecase{
		ch: ch,
	}

	go socket.consumeMessages()
	go socket.consumeChats()

	return socket
}

// consumeMessages принимает информацию о сообщениях (добавление/изменение/удаление)
func (w *WebsocketUsecase) consumeMessages() {
	for {
		messages, err := w.ch.Consume(
			"message", // queue
			"",        // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)

		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}
		for message := range messages {
			log.Printf("received a message: %s", message.Body)
		}

	}
}

// consumeChats принимает информацию об изменении чатов
func (w *WebsocketUsecase) consumeChats() {
	for {
		messages, err := w.ch.Consume(
			"chat", // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}
		for message := range messages {
			log.Printf("received a message: %s", message.Body)
		}
	}
}

func (w *WebsocketUsecase) addEventIntoChatRutine(event chatEvent.Event) {
	// если нет рутины чата, то сначала создадим ее
	if !w.isChatActive(event.ChatId) {
		
	}
}

func (w *WebsocketUsecase) isChatActive(chatId uuid.UUID) bool {
	_, ok := w.activeChats[chatId]
	return ok
}
