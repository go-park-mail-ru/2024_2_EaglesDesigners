package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	message "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	messageUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	admin = "admin"
	none  = "none"
	owner = "owner"
)

const (
	channel = "channel"
)

type ChatUsecaseImpl struct {
	tokenUsecase      *usecase.Usecase
	messageRepository message.MessageRepository
	repository        chatlist.ChatRepository
	activeUsers       map[uuid.UUID]bool
	redisClient       *redis.Client
}

func NewChatUsecase(tokenService *usecase.Usecase, repository chatlist.ChatRepository, messageRepository message.MessageRepository,
	activeUsers map[uuid.UUID]bool, redisClient *redis.Client) ChatUsecase {
	return &ChatUsecaseImpl{
		tokenUsecase:      tokenService,
		repository:        repository,
		messageRepository: messageRepository,
		activeUsers:       activeUsers,
		redisClient:       redisClient,
	}
}

func (s *ChatUsecaseImpl) createChatDTO(ctx context.Context, chat chatModel.Chat) (chatModel.ChatDTOOutput, error) {
	message, err := s.messageRepository.GetLastMessage(chat.ChatId)
	if err != nil {
		log.Printf("Usecase: не удалось получить последнее сообщение: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}
	log.Println("Usecase: последнее сообщение получено")

	log.Printf("Chat usecase: установка количества участников чата: %v", chat.ChatId)
	countOfUsers, err := s.repository.GetCountOfUsersInChat(ctx, chat.ChatId)

	if err != nil {
		log.Printf("Usecase: не удалось получить количество пользователей: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}
	log.Println("Usecase: количество пользователей получено")

	var photoBase64 string
	// Достаем фото из папки
	if chat.AvatarURL != "" {
		phId, err := uuid.Parse(chat.AvatarURL)
		if err != nil {
			return chatModel.ChatDTOOutput{}, err
		}

		photoBase64, err = base64helper.ReadPhotoBase64(phId)
		if err != nil && !os.IsNotExist(err) {
			return chatModel.ChatDTOOutput{}, err
		}
	}

	log.Println("Usecase: фото успешно считано и закодировано в base64")

	return chatModel.СhatToChatDTO(chat,
		countOfUsers,
		message,
		photoBase64), nil
}

func (s *ChatUsecaseImpl) GetChats(ctx context.Context, cookie []*http.Cookie, pageNum int) ([]chatModel.ChatDTOOutput, error) {

	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return []chatModel.ChatDTOOutput{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	log.Printf("Chat usecase: пришел запрос на получение всех чатов от пользователя: %v", user.ID)

	chats, err := s.repository.GetUserChats(ctx, user.ID, pageNum)
	if err != nil {
		return nil, err
	}
	log.Println("Usecase: чаты получены")

	chatsDTO := []chatModel.ChatDTOOutput{}

	for _, chat := range chats {

		chatDTO, err := s.createChatDTO(ctx, chat)

		if err != nil {
			log.Printf("Chat usecase -> GetChats: не удалось создать DTO: %v", err)
			return nil, err
		}

		chatsDTO = append(chatsDTO,
			chatDTO)
	}

	return chatsDTO, nil
}

func (s *ChatUsecaseImpl) sendNotificationToUser(ctx context.Context, userId uuid.UUID, chatDTO chatModel.ChatDTOOutput, method string) error {
	if _, ok := s.activeUsers[userId]; ok {

		log.Printf("Chat usecase -> sendNotificationToUser: начата отправка уведомления пользователю: %v", userId)
		err := s.redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: userId.String(),
			MaxLen: 0,
			ID:     "",
			Values: map[string]interface{}{
				method: chatDTO,
			},
		}).Err()

		if err != nil {
			log.Printf("Chat usecase -> sendNotificationToUser: уведомление не отправлено пользователю: %v. Ошибка: %v", userId, err)
		}

		return err
	}
	log.Printf("Chat usecase -> sendNotificationToUser: уведомление пользователю отправлено: %v", userId)
	return nil
}

func (s *ChatUsecaseImpl) AddUsersIntoChat(ctx context.Context, cookie []*http.Cookie, user_ids []uuid.UUID, chat_id uuid.UUID) error {
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	role, err := s.repository.GetUserRoleInChat(ctx, user.ID, chat_id)
	if err != nil {
		return err
	}

	// проверяем есть ли права
	switch role {
	case admin, owner:
		log.Printf("Chat usecase -> AddUsersIntoChat: начато добавление пользователей в чат %v пользователем %v", chat_id, user.ID)

		chat, err := s.repository.GetChatById(ctx, chat_id)

		if err != nil {
			log.Println("Chat usecase -> AddUsersIntoChat: не удалось добавить юзера в чат^ %v", err)
		}

		for _, id := range user_ids {
			s.repository.AddUserIntoChat(ctx, id, chat_id, none)
			chatDTO, err := s.createChatDTO(ctx, chat)
			if err != nil {
				log.Printf("Chat usecase -> AddUsersIntoChat: не удалось создать DTO: %v", err)
			}
			s.sendNotificationToUser(ctx, id, chatDTO, messageUsecase.FeatNewUser)

		}
		log.Printf("Chat usecase -> AddUsersIntoChat: участники добавлены в чат %v пользователем %v", chat_id, user.ID)

		return nil
	}

	return errors.New("Участники не добавлены")
}

func (s *ChatUsecaseImpl) AddNewChat(ctx context.Context, cookie []*http.Cookie, chat chatModel.ChatDTOInput) error {
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	//пока без фото
	// photoPath, err := base64helper.SavePhotoBase64(chat.AvatarBase64)

	// if err != nil {
	// 	log.Printf("Chat usecase -> AddNewChat: не удалось сохранить фото: %v", err)
	// 	return err
	// }

	chatId := uuid.New()

	newChat := chatModel.Chat{
		ChatId:    chatId,
		ChatName:  chat.ChatName,
		ChatType:  chat.ChatType,
	}

	if chat.ChatType == channel {
		newChat.ChatURLName = chat.ChatName + "_" + uuid.NewString()
	}

	// создание чата
	err = s.repository.CreateNewChat(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось сохнанить чат: %v", err)
		return err
	}

	// добавление владельца
	err = s.repository.AddUserIntoChat(ctx, user.ID, chatId, owner)

	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось добавить пользователя в чат: %v", err)
		return err
	}

	newChatDTO, err := s.createChatDTO(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось создать DTO: %v", err)
	}

	// добавляем основателя в чат
	s.sendNotificationToUser(ctx, user.ID, newChatDTO, messageUsecase.FeatNewUser)

	log.Printf("Chat usecase -> AddNewChat: начато добавление пользователей в чат. Количество бользователей на добавление: %v", len(chat.UsersToAdd))
	err = s.AddUsersIntoChat(ctx, cookie, chat.UsersToAdd, chatId)

	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось добавить пользователя в чат: %v", err)
		return err
	}

	return nil
}

func (s *ChatUsecaseImpl) DeleteChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) error {
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return err
	}

	// проверяем есть ли права
	switch role {
	case owner:
		log.Printf("Chat usecase -> DeleteChat: удаление чата %v", chatId)
		err = s.repository.DeleteChat(ctx, chatId)
		if err != nil {
			log.Printf("Chat usecase -> DeleteChat: не удалось удалить чат: %v", err)
			return err
		}
		return nil
	case none, admin:
		return &customerror.NoPermissionError{
			User: userId.String(),
			Area: fmt.Sprintf("чат: %v", chatId.String()),
		}
	}
	return nil
}
