package usecase

import (
	"context"
	"fmt"

	chatEvent "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"

	"github.com/google/uuid"
)

const Chat = "chat"
const Message = "message"

// consumeChats принимает информацию об изменении чатов
func (w *WebsocketUsecase) consumeChats() {
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
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
			event, err := chatEvent.DeserializeEvent(message.Body)
			if err != nil {
				log.Errorf("Невозморжно десериализовать оюъект: %v", err)
				continue
			}

			w.addChatEventIntoChatRutine(event)
		}
	}
}

func (w *WebsocketUsecase) addChatEventIntoChatRutine(event chatEvent.Event) {
	// если нет рутины чата, то сначала создадим ее
	if !w.isChatActive(event.ChatId) {
		err := w.initNewChatBroker(event.ChatId)
		if err != nil {
			return
		}
	}
	channel := w.onlineChats[event.ChatId].events
	channel <- event
}

func (w *WebsocketUsecase) isChatActive(chatId uuid.UUID) bool {
	_, ok := w.onlineChats[chatId]
	return ok
}

func (w *WebsocketUsecase) initNewChatBroker(chatId uuid.UUID) error {
	users, err := w.getOnlineUsersInChat(chatId)
	if err != nil {
		return err
	}
	w.onlineChats[chatId] = ChatInfo{
		events: make(chan chatEvent.Event, 10),
		users:  users,
	}
	go w.chatBroker(chatId)
	return nil
}

func (w *WebsocketUsecase) getOnlineUsersInChat(chatId uuid.UUID) (map[uuid.UUID]struct{}, error) {
	// ебашим в usecase Берем всех пользоватеолей
	users, err := w.chatRepository.GetUsersFromChat(context.Background(), chatId)
	if err != nil {
		log := logger.LoggerWithCtx(context.Background(), logger.Log)
		log.Errorf("Не удалось получить пользователей чата: %v", err)
		return nil, err
	}

	onlineUsersInChat := map[uuid.UUID]struct{}{}

	for _, user := range users {
		// если пользователь онлайн, то добавляем его в чат
		if _, ok := w.onlineUsers[user.ID]; ok {
			onlineUsersInChat[user.ID] = struct{}{}
		}
	}

	return onlineUsersInChat, nil
}

const (
	UpdateChat          = "updateChat"
	DeleteChat          = "deleteChat"
	NewChat             = "newChat"
	DeleteUsersFromChat = "delUsers"
	AddNewUsersInChat   = "addUsers"

	// пользователь стал онлайн
	AddWebcosketUser = "addWebSocketUser"
)

func (w *WebsocketUsecase) chatBroker(chatId uuid.UUID) {
	// здесь надо сходить посмотреть всех юзеров
	chatInfo := w.onlineChats[chatId]
	events := chatInfo.events
	users := chatInfo.users
	for {
		// если ноль пользователей онлайн, то закрываем брокер
		if len(users) == 0 {
			log := logger.LoggerWithCtx(context.Background(), logger.Log)
			log.Debugf("Брокер для чата %v закрывается", chatId)

			close(w.onlineChats[chatId].events)
			delete(w.onlineChats, chatId)
			return
		}

		newEvent := <-events
		fmt.Println(newEvent)
		switch newEvent.Action {
		case AddWebcosketUser:
			// достаем из массива event.Users пользователя
			for _, userId := range newEvent.Users {
				users[userId] = struct{}{}
			}
		case DeleteChat:
			go w.sendEventToAllUsers(users, newEvent)
			delete(w.onlineChats, chatId)
			return
		case NewChat, UpdateChat:
			go w.sendEventToAllUsers(users, newEvent)
		case DeleteUsersFromChat:
			w.sendEventToDeletedUsers(newEvent.Users, newEvent)
			go w.sendEventToAllUsers(users, newEvent)
			// удаляем юзеров, если они были в подписчиках
			for _, userId := range newEvent.Users {
				if _, ok := w.onlineUsers[userId]; ok {
					delete(users, userId)
				}
			}
		case AddNewUsersInChat:
			// если пользователь онлайны, то добавляем в подписчики
			for _, userId := range newEvent.Users {
				if _, ok := w.onlineUsers[userId]; ok {
					users[userId] = struct{}{}
				}
			}
			go w.sendEventToAllUsers(users, newEvent)
		}
	}
}

func (w *WebsocketUsecase) sendEventToAllUsers(users map[uuid.UUID]struct{}, event chatEvent.Event) {
	for userId := range users {
		if _, ok := w.onlineUsers[userId]; ok {
			w.onlineUsers[userId] <- AnyEvent{
				TypeOfEvent: Chat,
				Event:       event,
			}
		}
	}
}

const (
	// current - имеется ввиду пользователь, который щас подписан на вебсокет
	CurrentUserDeleeted = "currentUserDeleeted"
)

func (w *WebsocketUsecase) sendEventToDeletedUsers(users []uuid.UUID, event chatEvent.Event) {
	event.Action = CurrentUserDeleeted
	event.Users = nil

	for _, userId := range users {
		if _, ok := w.onlineUsers[userId]; ok {
			w.onlineUsers[userId] <- AnyEvent{
				TypeOfEvent: Chat,
				Event:       event,
			}
		}
	}
}
