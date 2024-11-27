package usecase

import (
	"context"
	"encoding/json"

	messageModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
)

type MessageEvent struct {
	Action  string `json:"action"`
	Message messageModel.Message `json:"payload"`
}

const (
	DeleteMessage = "deleteMessage"
	NewMessage    = "newMessage"
	UpdateMessage = "updateMessage"
)

func SerializeMessageEvent(event MessageEvent) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeMessageEvent(data []byte) (MessageEvent, error) {
	var event MessageEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return MessageEvent{}, err
	}
	return event, nil
}

// consumeMessages принимает информацию о сообщениях (добавление/изменение/удаление)
func (w *WebsocketUsecase) consumeMessages() {
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
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
			msg, err := DeserializeMessageEvent(message.Body)

			if err != nil {
				log.Errorf("Невозморжно десериализовать оюъект: %v", err)
				continue
			}
			if _, ok := w.onlineChats[msg.Message.ChatId]; !ok {
				w.initNewChatBroker(msg.Message.ChatId)
			}
			w.sendMessage(msg)
		}

	}
}

func (w *WebsocketUsecase) sendMessage(event MessageEvent) {
	chatId := event.Message.ChatId
	users := w.onlineChats[chatId].users

	for user := range users {
		w.onlineUsers[user] <- AnyEvent{
			TypeOfEvent: Message,
			Event:       event,
		}
	}
}
