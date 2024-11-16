package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/custom_error"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	message "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	multipartHepler "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/multipartHelper"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/validator"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
	errGroup "golang.org/x/sync/errgroup"
)

const (
	admin     = "admin"
	none      = "none"
	owner     = "owner"
	notInChat = ""
)
const (
	personal = "personal"
	channel  = "channel"
)

const chatDir = "chat"

// ивенты для сокета
const (
	UpdateChat          = "updateChat"
	DeleteChat          = "deleteChat"
	NewChat             = "newChat"
	DeleteUsersFromChat = "delUsers"
	AddNewUsersInChat   = "addUsers"
)

type ChatUsecaseImpl struct {
	tokenUsecase      *usecase.Usecase
	messageRepository message.MessageRepository
	repository        chatlist.ChatRepository
	chatQuery         string
	ch                *amqp.Channel
}

func NewChatUsecase(tokenService *usecase.Usecase, repository chatlist.ChatRepository, messageRepository message.MessageRepository, ch *amqp.Channel) ChatUsecase {
	// объявляем очередь для яатов
	q, err := ch.QueueDeclare(
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

	return &ChatUsecaseImpl{
		tokenUsecase:      tokenService,
		repository:        repository,
		messageRepository: messageRepository,
		chatQuery:         q.Name,
		ch:                ch,
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

func (s *ChatUsecaseImpl) AddUsersIntoChat(ctx context.Context, cookie []*http.Cookie, user_ids []uuid.UUID, chatId uuid.UUID) (chatModel.AddedUsersIntoChatDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return chatModel.AddedUsersIntoChatDTO{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	role, err := s.repository.GetUserRoleInChat(ctx, user.ID, chatId)
	if err != nil {
		return chatModel.AddedUsersIntoChatDTO{}, err
	}

	var addedUsers []uuid.UUID
	var notAddedUsers []uuid.UUID
	// проверяем есть ли права
	switch role {
	case admin, owner, none:
		log.Printf("Chat usecase -> AddUsersIntoChat: начато добавление пользователей в чат %v пользователем %v", chatId, user.ID)

		for _, id := range user_ids {
			err = s.repository.AddUserIntoChat(ctx, id, chatId, none)
			if err != nil {
				notAddedUsers = append(notAddedUsers, id)
				continue
			}
			addedUsers = append(addedUsers, id)
		}
		log.Printf("Chat usecase -> AddUsersIntoChat: участники добавлены в чат %v пользователем %v", chatId, user.ID)

		s.sendIvent(ctx, AddNewUsersInChat, chatId, addedUsers)
		return chatModel.AddedUsersIntoChatDTO{AddedUsers: addedUsers,
			NotAddedUsers: notAddedUsers}, nil
	default:
		return chatModel.AddedUsersIntoChatDTO{}, &customerror.NoPermissionError{
			User: user.ID.String(),
			Area: fmt.Sprintf("Чат %v", chatId),
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

	// отправляем уведомление
	s.sendIvent(ctx, NewChat, chatId, nil)

	// добавляем основателя в чат

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
	for _, u := range users {
		if u.ID != user.ID {
			// находим имя пользователя и аватар
			chatName, avatar, err := s.repository.GetNameAndAvatar(ctx, u.ID)

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

		s.sendIvent(ctx, DeleteChat, chatId, nil)
		return nil
	case none, admin:
		log.Printf("У пользователя %v нет прав на удаление чата %v", userId, chatId)
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

		// кидаем уведомление в сокет
		s.sendIvent(ctx, UpdateChat, chatId, nil)
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

		// кидаем в веб сокет
		s.sendIvent(ctx, DeleteUsersFromChat, chatId, deletedIds)
		return chatModel.DeletdeUsersFromChatDTO{DeletedUsers: deletedIds}, nil

	default:
		return chatModel.DeletdeUsersFromChatDTO{}, errors.New("Участники не удалены")
	}
}

func (s *ChatUsecaseImpl) GetUsersFromChat(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (chatModel.UsersInChatDTO, error) {
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return chatModel.UsersInChatDTO{}, err
	}

	if role == notInChat {
		return chatModel.UsersInChatDTO{}, &customerror.NoPermissionError{
			User: userId.String(),
			Area: chatId.String(),
		}
	}

	users, err := s.repository.GetUsersFromChat(ctx, chatId)
	if err != nil {
		return chatModel.UsersInChatDTO{}, err
	}

	return convertUsersInChatToDTO(users), nil
}

func convertUsersInChatToDTO(users []chatModel.UserInChatDAO) chatModel.UsersInChatDTO {
	var usersDTO []chatModel.UserInChatDTO

	var mu sync.Mutex
	var g errGroup.Group

	for _, user := range users {
		g.Go(func() error {
			userDTO := chatModel.UserInChatDTO{
				ID:         user.ID,
				Username:   html.EscapeString(user.Username),
				Name:       validator.EscapePtrString(user.Name),
				AvatarPath: validator.EscapePtrString(user.AvatarPath),
			}

			if user.Role != nil {
				userDTO.Role = new(string)
				switch *user.Role {
				case 1:
					*userDTO.Role = none
				case 2:
					*userDTO.Role = owner
				case 3:
					*userDTO.Role = admin
				}
			}

			mu.Lock()
			defer mu.Unlock()

			usersDTO = append(usersDTO, userDTO)

			return nil
		})
	}

	g.Wait()

	return chatModel.UsersInChatDTO{
		Users: usersDTO,
	}
}

func (s *ChatUsecaseImpl) sendIvent(ctx context.Context, action string, chatId uuid.UUID, users []uuid.UUID) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	newEvent := chatModel.Event{
		Action: action,
		ChatId: chatId,
		Users:  users,
	}

	body, err := chatModel.SerializeEvent(newEvent)
	if err != nil {
		log.Errorf("Не удалось сериализовать объект")
		return
	}
	err = s.ch.PublishWithContext(ctx,
		"",          // exchange
		s.chatQuery, // имя очереди
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
