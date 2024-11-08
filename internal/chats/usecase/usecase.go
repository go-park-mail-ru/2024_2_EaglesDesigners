package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	message "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	messageUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/multipartHelper"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	admin     = "admin"
	none      = "none"
	owner     = "owner"
	notInChat = ""
)
const (
	personal = "personal"
)

const (
	channel = "channel"
)

const chatDir = "chat"

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
	log := logger.LoggerWithCtx(ctx, logger.Log)
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

	return chatModel.СhatToChatDTO(chat,
		countOfUsers,
		message), nil
}

func (s *ChatUsecaseImpl) GetChats(ctx context.Context, cookie []*http.Cookie, pageNum int) ([]chatModel.ChatDTOOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
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
		if chat.ChatType == personal {
			chat.ChatName, chat.AvatarURL, err = s.getAvatarAndNameForPersonalChat(ctx, user, chat.ChatId)

			if err != nil {
				log.Errorf("Chat usecase -> GetChats: не удалось обработать персональный чат: %v", err)
				return nil, err
			}
		}

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
	log := logger.LoggerWithCtx(ctx, logger.Log)
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

func (s *ChatUsecaseImpl) AddUsersIntoChat(ctx context.Context, cookie []*http.Cookie, user_ids []uuid.UUID, chat_id uuid.UUID) (chatModel.AddedUsersIntoChatDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return chatModel.AddedUsersIntoChatDTO{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	role, err := s.repository.GetUserRoleInChat(ctx, user.ID, chat_id)
	if err != nil {
		return chatModel.AddedUsersIntoChatDTO{}, err
	}

	var addedUsers []uuid.UUID
	var notAddedUsers []uuid.UUID
	// проверяем есть ли права
	switch role {
	case admin, owner, none:
		log.Printf("Chat usecase -> AddUsersIntoChat: начато добавление пользователей в чат %v пользователем %v", chat_id, user.ID)

		chat, err := s.repository.GetChatById(ctx, chat_id)

		if err != nil {
			log.Println("Chat usecase -> AddUsersIntoChat: не удалось добавить юзера в чат^ %v", err)
		}

		for _, id := range user_ids {
			err = s.repository.AddUserIntoChat(ctx, id, chat_id, none)
			if err != nil {
				notAddedUsers = append(notAddedUsers, id)
				continue
			}
			chatDTO, err := s.createChatDTO(ctx, chat)
			if err != nil {
				log.Printf("Chat usecase -> AddUsersIntoChat: не удалось создать DTO: %v", err)
			}
			s.sendNotificationToUser(ctx, id, chatDTO, messageUsecase.FeatNewUser)
			addedUsers = append(addedUsers, id)
		}
		log.Printf("Chat usecase -> AddUsersIntoChat: участники добавлены в чат %v пользователем %v", chat_id, user.ID)

		return chatModel.AddedUsersIntoChatDTO{AddedUsers: addedUsers,
			NotAddedUsers: notAddedUsers}, nil
	default:
		return chatModel.AddedUsersIntoChatDTO{}, &customerror.NoPermissionError{
			User: user.ID.String(),
			Area: fmt.Sprintf("Чат %v", chat_id),
		}
	}
}

func (s *ChatUsecaseImpl) AddNewChat(ctx context.Context, cookie []*http.Cookie, chat chatModel.ChatDTOInput) (chatModel.ChatDTOOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return chatModel.ChatDTOOutput{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	chatId := uuid.New()

	newChat := chatModel.Chat{
		ChatId:   chatId,
		ChatName: chat.ChatName,
		ChatType: chat.ChatType,
	}

	if chat.Avatar != nil {
		filename, err := multipartHepler.SavePhoto(*chat.Avatar, chatDir)
		if err != nil {
			log.Printf("Не удалось записать аватарку: %v", err)
			return chatModel.ChatDTOOutput{}, err
		}
		newChat.AvatarURL = filename
	}

	if chat.ChatType == channel {
		newChat.ChatURLName = chat.ChatName + "_" + uuid.NewString()
	}

	// создание чата
	err = s.repository.CreateNewChat(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось сохнанить чат: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}

	// добавление владельца
	err = s.repository.AddUserIntoChat(ctx, user.ID, chatId, owner)

	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось добавить пользователя в чат: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}

	newChatDTO, err := s.createChatDTO(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось создать DTO: %v", err)
	}

	// добавляем основателя в чат
	s.sendNotificationToUser(ctx, user.ID, newChatDTO, messageUsecase.FeatNewUser)

	log.Printf("Chat usecase -> AddNewChat: начато добавление пользователей в чат. Количество бользователей на добавление: %v", len(chat.UsersToAdd))
	_, err = s.AddUsersIntoChat(ctx, cookie, chat.UsersToAdd, chatId)

	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось добавить пользователя в чат: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}

	if newChatDTO.ChatType == personal {
		newChatDTO.ChatName, newChatDTO.AvatarPath, err = s.getAvatarAndNameForPersonalChat(ctx, user, newChat.ChatId)

		if err != nil {
			log.Errorf("Chat usecase -> AddNewChat: не удалось обработать персональный чат: %v", err)
			return chatModel.ChatDTOOutput{}, err
		}
	}

	return newChatDTO, nil
}

func (s *ChatUsecaseImpl) getAvatarAndNameForPersonalChat(ctx context.Context, user usecase.User, chatId uuid.UUID) (string, string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	users, err := s.repository.GetUsersFromChat(ctx, chatId)
	if err != nil {
		return "", "", err
	}
	for _, id := range users {
		if id != user.ID {
			// находим имя пользователя и аватар
			chatName, avatar, err := s.repository.GetUserNameAndAvatar(ctx, id)

			if err != nil {
				log.Printf("Chat usecase -> GetChats: не удалось получить аватар и имя: %v", err)
				return "", "", err
			}
			return chatName, avatar, err
		}
	}
	return "", "", nil
}

func (s *ChatUsecaseImpl) DeleteChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return err
	}

	// проверяем есть ли права
	switch role {
	case owner:
		log.Printf("Chat usecase -> DeleteChat: удаление чата %v", chatId)

		// send notification to chat

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

func (s *ChatUsecaseImpl) UpdateChat(ctx context.Context, chatId uuid.UUID, chatUpdate chatModel.ChatUpdate, userId uuid.UUID) (chatModel.ChatUpdateOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return chatModel.ChatUpdateOutput{}, err
	}

	var updatedChat chatModel.ChatUpdateOutput
	// проверяем есть ли права
	switch role {
	case owner, admin, none:
		log.Printf("Chat usecase -> UpdateChat: обновление чата %v", chatId)

		// send notification to chat
		if chatUpdate.Avatar != nil {
			chat, err := s.repository.GetChatById(ctx, chatId)
			if err != nil {
				return chatModel.ChatUpdateOutput{}, err
			}

			if chat.AvatarURL != "" {

				err = multipartHepler.RewritePhoto(*chatUpdate.Avatar, chat.AvatarURL)
				if err != nil {
					log.Printf("Chat usecase -> UpdateChat: не удалось обновить аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				updatedChat.Avatar = chat.AvatarURL
			} else {
				log.Println("Chat usecase -> UpdateChat: нет старой аватарки -> установка новой")
				filename, err := multipartHepler.SavePhoto(*chatUpdate.Avatar, chatDir)
				if err != nil {
					log.Printf("Не удалось записать аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				err = s.repository.UpdateChatPhoto(ctx, chatId, filename)

				if err != nil {
					log.Printf("Chat usecase -> UpdateChat: не удалось установить аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				updatedChat.Avatar = filename
			}
			log.Println("Chat usecase -> UpdateChat: аватар обновлен")

		}

		if chatUpdate.ChatName != "" {
			err = s.repository.UpdateChat(ctx, chatId, chatUpdate.ChatName)
			if err != nil {
				log.Printf("Chat usecase -> UpdateChat: не удалось обновить имя чата: %v", err)
				return chatModel.ChatUpdateOutput{}, err
			}
			log.Println("Chat usecase -> UpdateChat: имя чата обновлено")
			updatedChat.ChatName = chatUpdate.ChatName
		}
		return updatedChat, nil
	default:
		log.Printf("Chat usecase -> UpdateChat: у пользователя %v нет привелегий", userId)
		return chatModel.ChatUpdateOutput{}, &customerror.NoPermissionError{
			User: userId.String(),
			Area: fmt.Sprintf("чат: %v", chatId.String()),
		}
	}
}

func (s *ChatUsecaseImpl) DeleteUsersFromChat(ctx context.Context, userID uuid.UUID, chatId uuid.UUID, usertToDelete chatModel.DeleteUsersFromChatDTO) (chatModel.DeletdeUsersFromChatDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	role, err := s.repository.GetUserRoleInChat(ctx, userID, chatId)
	if err != nil {
		return chatModel.DeletdeUsersFromChatDTO{}, err
	}
	var deletedIds []uuid.UUID
	// проверяем есть ли права
	switch role {
	case admin, owner:
		log.Printf("Chat usecase -> DeleteUsersFromChat: начато удаление пользователей в чат %v пользователем %v", chatId, userID)

		for _, id := range usertToDelete.UsersId {
			userRole, err := s.repository.GetUserRoleInChat(ctx, id, chatId)
			if err != nil {
				continue
			}

			if id == userID {
				continue
			}
			if userRole == owner {
				continue
			}

			err = s.repository.DeleteUserFromChat(ctx, id, chatId)
			if err != nil {
				continue
			}
			deletedIds = append(deletedIds, id)
		}
		log.Printf("Chat usecase -> DeleteUsersFromChat: участники удалены из чата %v пользователем %v", chatId, userID)

		return chatModel.DeletdeUsersFromChatDTO{}, nil
	}

	return chatModel.DeletdeUsersFromChatDTO{DeletedUsers: deletedIds}, errors.New("Участники не удалены")
}

func (s *ChatUsecaseImpl) GetUsersFromChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (chatModel.UsersInChat, error) {
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return chatModel.UsersInChat{}, err
	}

	if role == notInChat {
		return chatModel.UsersInChat{}, &customerror.NoPermissionError{
			User: userId.String(),
			Area: chatId.String(),
		}
	}

	ids, err := s.repository.GetUsersFromChat(ctx, chatId)
	if err != nil {
		return chatModel.UsersInChat{}, err
	}

	return chatModel.UsersInChat{
		UsersId: ids,
	}, nil
}
