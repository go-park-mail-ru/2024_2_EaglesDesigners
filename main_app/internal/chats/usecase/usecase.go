package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
	customerror "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/custom_error"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/repository"
	messageModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
	message "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/validator"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
	errGroup "golang.org/x/sync/errgroup"
)

func init() {
	prometheus.MustRegister(addNewChatMetric, addUsersIntoChatMetric)
}

const (
	Admin     = "admin"
	None      = "none"
	Owner     = "owner"
	NotInChat = ""
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
	fileUC         FilesUsecase
	messageUsecase message.MessageUsecase
	repository     chatlist.ChatRepository
	chatQuery      string
	ch             *amqp.Channel
}

type FilesUsecase interface {
	RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) error
	DeletePhoto(ctx context.Context, fileIDStr string) error
	IsImage(header multipart.FileHeader) error
	SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}

func NewChatUsecase(fileUC FilesUsecase, repository chatlist.ChatRepository, messageRepository message.MessageUsecase, ch *amqp.Channel) ChatUsecase {
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
		fileUC:         fileUC,
		messageUsecase: messageRepository,
		repository:     repository,
		chatQuery:      q.Name,
		ch:             ch,
	}
}

func (s *ChatUsecaseImpl) createChatDTO(ctx context.Context, chat chatModel.Chat) (chatModel.ChatDTOOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	message, err := s.messageUsecase.GetLastMessage(chat.ChatId)
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

func (s *ChatUsecaseImpl) GetChats(ctx context.Context, cookie []*http.Cookie) ([]chatModel.ChatDTOOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		return nil, errors.New(responser.UserNotFoundError)
	}
	log.Printf("Chat usecase: пришел запрос на получение всех чатов от пользователя: %v", user.ID)

	chats, err := s.repository.GetUserChats(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	log.Println("Usecase: чаты получены")

	chatsDTO := []chatModel.ChatDTOOutput{}

	for _, chat := range chats {
		if chat.ChatType == personal {
			chat.ChatName, chat.AvatarURL, err = s.getAvatarAndNameForPersonalChat(ctx, user.ID, chat.ChatId)

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

	sort.Sort(chatModel.ByLastMessage(chatsDTO))

	return chatsDTO, nil
}

var addUsersIntoChatMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_added_users_into_chat",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

func (s *ChatUsecaseImpl) addUsersIntoChat(ctx context.Context, user_ids []uuid.UUID, chatId uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
	var addedUsers []uuid.UUID
	var notAddedUsers []uuid.UUID
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("начато добавление пользователей в чат %v", chatId)

	for _, id := range user_ids {
		err := s.repository.AddUserIntoChat(ctx, id, chatId, None)
		if err != nil {
			notAddedUsers = append(notAddedUsers, id)
			continue
		}
		addedUsers = append(addedUsers, id)
	}
	log.Printf("Участники добавлены в чат %v", chatId)
	metric.IncMetric(*addUsersIntoChatMetric)
	return addedUsers, notAddedUsers
}

func (s *ChatUsecaseImpl) AddUsersIntoChatWithCheckPermission(ctx context.Context, userIds []uuid.UUID, chatId uuid.UUID) (chatModel.AddedUsersIntoChatDTO, error) {
	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
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
	case Admin, Owner, None:
		addedUsers, notAddedUsers = s.addUsersIntoChat(ctx, userIds, chatId)
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

var addNewChatMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_added_chats",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

func (s *ChatUsecaseImpl) AddNewChat(ctx context.Context, cookie []*http.Cookie, chat chatModel.ChatDTOInput) (chatModel.ChatDTOOutput, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	user, ok := ctx.Value(auth.UserKey).(auth.User)
	if !ok {
		return chatModel.ChatDTOOutput{}, errors.New(responser.UserNotFoundError)
	}

	chatId := uuid.New()

	newChat := chatModel.Chat{
		ChatId:   chatId,
		ChatName: chat.ChatName,
		ChatType: chat.ChatType,
	}

	if chat.Avatar != nil {
		filename, err := s.fileUC.SaveAvatar(ctx, *chat.Avatar, chat.AvatarHeader)
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
	err := s.repository.CreateNewChat(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось сохнанить чат: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}

	// добавление владельца
	err = s.repository.AddUserIntoChat(ctx, user.ID, chatId, Owner)

	if err != nil {
		log.Printf("Не удалось добавить владельца в чат: %v", err)
		return chatModel.ChatDTOOutput{}, err
	}

	newChatDTO, err := s.createChatDTO(ctx, newChat)
	if err != nil {
		log.Printf("Chat usecase -> AddNewChat: не удалось создать DTO: %v", err)
	}

	// добавляем пользователей в чат
	log.Printf("Chat usecase -> AddNewChat: начато добавление пользователей в чат. Количество бользователей на добавление: %v", len(chat.UsersToAdd))
	s.addUsersIntoChat(ctx, chat.UsersToAdd, chatId)

	if newChatDTO.ChatType == personal {
		newChatDTO.ChatName, newChatDTO.AvatarPath, err = s.getAvatarAndNameForPersonalChat(ctx, user.ID, newChat.ChatId)

		if err != nil {
			log.Errorf("Chat usecase -> AddNewChat: не удалось обработать персональный чат: %v", err)
			return chatModel.ChatDTOOutput{}, err
		}
	}

	// отправляем уведомлениея
	s.sendIvent(ctx, NewChat, chatId, nil)

	createChatMessage := "Группа создана."
	switch chat.ChatType {
	case channel:
		createChatMessage = "Канал создан."
	case personal:
		createChatMessage = "Личный чат создан."
	}

	newChatDTO.LastMessage = s.sendInformationalMessage(ctx, user.ID, chatId, createChatMessage)
	metric.IncMetric(*addNewChatMetric)
	return newChatDTO, nil
}

func (s *ChatUsecaseImpl) sendInformationalMessage(ctx context.Context, userID uuid.UUID, chatId uuid.UUID, event string) messageModel.Message {
	message := messageModel.Message{
		MessageId:   uuid.New(),
		AuthorID:    userID, // При выдаче, если информационное, будет меняться uuid юзера на нули.
		Message:     event,
		SentAt:      time.Now(),
		MessageType: "informational",
	}
	s.messageUsecase.SendInformationalMessage(ctx, message, chatId)
	message.AuthorID, _ = uuid.Parse("00000000-0000-0000-0000-000000000000") // Для информационных сообщений ставим нули в автора.
	return message
}

func (s *ChatUsecaseImpl) getAvatarAndNameForPersonalChat(ctx context.Context, userID uuid.UUID, chatId uuid.UUID) (string, string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	users, err := s.repository.GetUsersFromChat(ctx, chatId)
	if err != nil {
		return "", "", err
	}
	for _, u := range users {
		if u.ID != userID {
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

	chatType, err := s.repository.GetChatType(ctx, chatId)
	if err != nil {
		return err
	}

	if chatType == personal {
		log.Printf("Chat usecase -> DeleteChat: удаление чата %v", chatId)

		// send notification to chat

		err = s.repository.DeleteChat(ctx, chatId)
		if err != nil {
			log.Printf("Chat usecase -> DeleteChat: не удалось удалить чат: %v", err)
			return err
		}

		s.sendIvent(ctx, DeleteChat, chatId, nil)
		return nil
	}

	// проверяем есть ли права
	switch role {
	case Owner:
		log.Printf("Chat usecase -> DeleteChat: удаление чата %v", chatId)

		// send notification to chat

		err = s.repository.DeleteChat(ctx, chatId)
		if err != nil {
			log.Printf("Chat usecase -> DeleteChat: не удалось удалить чат: %v", err)
			return err
		}

		s.sendIvent(ctx, DeleteChat, chatId, nil)
		return nil
	case None, Admin:
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

	chatType, err := s.repository.GetChatType(ctx, chatId)
	if err != nil {
		return chatModel.ChatUpdateOutput{}, err
	}

	// по умолчанию равно false
	var hasPermission bool

	// проверяем есть ли права
	if chatType == channel {
		switch role {
		case Owner, Admin:
			hasPermission = true
		}
	} else {
		switch role {
		case Owner, Admin, None:
			hasPermission = true
		}
	}

	var updatedChat chatModel.ChatUpdateOutput

	if hasPermission {
		log.Printf("обновление чата %v", chatId)

		// send notification to chat
		if chatUpdate.Avatar != nil {
			chat, err := s.repository.GetChatById(ctx, chatId)
			if err != nil {
				return chatModel.ChatUpdateOutput{}, err
			}

			if chat.AvatarURL != "" {

				err = s.fileUC.RewritePhoto(ctx, *chatUpdate.Avatar, *chatUpdate.AvatarHeader, chat.AvatarURL)
				if err != nil {
					log.Errorf("не удалось обновить аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				updatedChat.Avatar = chat.AvatarURL
			} else {
				log.Println("нет старой аватарки -> установка новой")
				filename, err := s.fileUC.SaveAvatar(ctx, *chatUpdate.Avatar, chatUpdate.AvatarHeader)
				if err != nil {
					log.Errorf("Не удалось записать аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				err = s.repository.UpdateChatPhoto(ctx, chatId, filename)

				if err != nil {
					log.Errorf("не удалось установить аватарку: %v", err)
					return chatModel.ChatUpdateOutput{}, err
				}
				updatedChat.Avatar = filename
			}
			log.Println("аватар обновлен")

		}

		if chatUpdate.ChatName != "" {
			err = s.repository.UpdateChat(ctx, chatId, chatUpdate.ChatName)
			if err != nil {
				log.Errorf("не удалось обновить имя чата: %v", err)
				return chatModel.ChatUpdateOutput{}, err
			}
			log.Println("имя чата обновлено")
			updatedChat.ChatName = chatUpdate.ChatName
		}

		// кидаем уведомление в сокет
		s.sendIvent(ctx, UpdateChat, chatId, nil)
		return updatedChat, nil
	} else {
		log.Printf("у пользователя %v нет привелегий", userId)
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
	case Admin, Owner:
		log.Printf("Chat usecase -> DeleteUsersFromChat: начато удаление пользователей в чат %v пользователем %v", chatId, userID)

		for _, id := range usertToDelete.UsersId {
			userRole, err := s.repository.GetUserRoleInChat(ctx, id, chatId)
			if err != nil {
				continue
			}

			if id == userID {
				continue
			}
			if userRole == Owner {
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
		return chatModel.DeletdeUsersFromChatDTO{}, errors.New("участники не удалены")
	}
}

// UserLeaveChat удаляет владельца обращения из чата
func (s *ChatUsecaseImpl) UserLeaveChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return err
	}
	if role == NotInChat {
		log.Printf("Пользователь %v не состоит в чате %v", userId, chatId)
		return &customerror.NoPermissionError{
			User: userId.String(),
			Area: fmt.Sprintf("Пользователь %v не состоит в чате %v", userId, chatId),
		}
	}

	if role == Owner {
		return &customerror.NoPermissionError{
			User: userId.String(),
			Area: fmt.Sprintf("Пользователь %v является владельцем чата %v", userId, chatId),
		}
	}

	err = s.repository.DeleteUserFromChat(ctx, userId, chatId)
	if err != nil {
		log.Printf("Не удалось удалить пользователя %v из чата %v", userId, chatId)
		return err
	}

	s.sendIvent(ctx, DeleteUsersFromChat, chatId, []uuid.UUID{userId})
	return nil
}

func (s *ChatUsecaseImpl) GetChatInfo(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (chatModel.ChatInfoDTO, error) {
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return chatModel.ChatInfoDTO{}, err
	}

	if role == NotInChat {
		return chatModel.ChatInfoDTO{}, &customerror.NoPermissionError{
			User: userId.String(),
			Area: chatId.String(),
		}
	}

	var g errGroup.Group

	var users []chatModel.UserInChatDAO
	var usersDTO []chatModel.UserInChatDTO
	var messages []messageModel.Message
	var files []messageModel.Payload
	var photos []messageModel.Payload

	g.Go(func() error {
		users, err = s.repository.GetUsersFromChat(ctx, chatId)
		if err != nil {
			return err
		}
		usersDTO = convertUsersInChatToDTO(users)

		return nil
	})

	g.Go(func() error {
		messagesDTO, err := s.messageUsecase.GetFirstMessages(ctx, chatId)
		messages = messagesDTO.Messages
		return err
	})

	g.Go(func() error {
		files, photos, err = s.messageUsecase.GetPayload(ctx, chatId)
		return err
	})

	if err := g.Wait(); err != nil {
		return chatModel.ChatInfoDTO{}, err
	}

	sendNotifications, err := s.repository.GetSendNotificationsForUser(ctx, chatId, userId)

	return chatModel.ChatInfoDTO{
		Role:              role,
		Users:             usersDTO,
		Messages:          messages,
		SendNotifications: sendNotifications,
		Files:             files,
		Photos:            photos,
	}, nil
}

func (s *ChatUsecaseImpl) AddBranch(ctx context.Context, chatId uuid.UUID, messageID uuid.UUID, userId uuid.UUID) (chatModel.AddBranch, error) {
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatId)
	if err != nil {
		return chatModel.AddBranch{}, err
	}

	if role == NotInChat {
		return chatModel.AddBranch{}, &customerror.NoPermissionError{
			User: userId.String(),
			Area: chatId.String(),
		}
	}

	chatType, err := s.repository.GetChatType(ctx, chatId)
	if err != nil {
		return chatModel.AddBranch{}, err
	}

	if chatType == personal {
		return chatModel.AddBranch{}, errors.New("нельзя добавить ветку в личный чат")
	}

	branch, err := s.repository.AddBranch(ctx, chatId, messageID)
	if err != nil {
		return chatModel.AddBranch{}, err
	}

	return branch, nil
}

func (s *ChatUsecaseImpl) SearchChats(ctx context.Context, userID uuid.UUID, keyWord string) (chatModel.SearchChatsDTO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Debugf("пришел запрос на получение всех чатов от пользователя: %v", userID)

	var g errGroup.Group

	var userChatsDTO []chatModel.ChatDTOOutput
	var globalChannelsDTO []chatModel.ChatDTOOutput

	g.Go(func() error {
		userChats, err := s.repository.SearchUserChats(ctx, userID, keyWord)
		if err != nil {
			return err
		}
		log.Debugln("чаты пользователя получены")

		for _, chat := range userChats {
			if chat.ChatType == personal {
				chat.ChatName, chat.AvatarURL, err = s.getAvatarAndNameForPersonalChat(ctx, userID, chat.ChatId)

				if err != nil {
					log.Errorf("не удалось обработать персональный чат: %v", err)
					return err
				}
			}

			chatDTO, err := s.createChatDTO(ctx, chat)

			if err != nil {
				log.Errorf("не удалось создать DTO: %v", err)
				return err
			}

			userChatsDTO = append(userChatsDTO,
				chatDTO)
		}

		return nil
	})

	g.Go(func() error {
		globalChannels, err := s.repository.SearchGlobalChats(ctx, userID, keyWord)
		if err != nil {
			return err
		}
		log.Debugln("глобальные каналы получены")

		for _, chat := range globalChannels {
			channelDTO, err := s.createChatDTO(ctx, chat)

			if err != nil {
				log.Errorf("не удалось создать DTO: %v", err)
				return err
			}

			globalChannelsDTO = append(globalChannelsDTO,
				channelDTO)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return chatModel.SearchChatsDTO{}, err
	}

	outputDTO := chatModel.SearchChatsDTO{
		UserChats:      userChatsDTO,
		GlobalChannels: globalChannelsDTO,
	}

	return outputDTO, nil
}

func (s *ChatUsecaseImpl) JoinChannel(ctx context.Context, userId uuid.UUID, channelId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	chatType, err := s.repository.GetChatType(ctx, channelId)
	if err != nil {
		return err
	}
	if chatType != channel {
		return &customerror.NoPermissionError{
			Area: fmt.Sprintf("%v не явяляется каналом", channelId),
		}
	}
	_, notAdded := s.addUsersIntoChat(ctx, []uuid.UUID{userId}, channelId)
	if len(notAdded) != 0 {
		log.Errorf("Пользователю %v не удалось вступить в канал %v", userId, channelId)
		return errors.New("не удалось добавить пользователя в чат")
	}
	return nil
}

func convertUsersInChatToDTO(users []chatModel.UserInChatDAO) []chatModel.UserInChatDTO {
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
					*userDTO.Role = None
				case 2:
					*userDTO.Role = Owner
				case 3:
					*userDTO.Role = Admin
				}
			}

			mu.Lock()
			defer mu.Unlock()

			usersDTO = append(usersDTO, userDTO)

			return nil
		})
	}

	g.Wait()

	return usersDTO
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

func (s *ChatUsecaseImpl) GetUserChats(ctx context.Context, userId string) (chatIds []string, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	chatIds = make([]string, 0)
	chats, err := s.repository.GetUserChats(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	for _, chat := range chats {
		chatIds = append(chatIds, chat.ChatId.String())
	}

	return chatIds, nil
}

func (s *ChatUsecaseImpl) GetUsersFromChat(ctx context.Context, chatId string) (userIds []string, err error) {
	chatUUID, err := uuid.Parse(chatId)
	if err != nil {
		return nil, err
	}

	userIds = make([]string, 0)
	users, err := s.repository.GetUsersFromChat(ctx, chatUUID)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userIds = append(userIds, user.ID.String())
	}
	return userIds, nil
}

// SetChatNotofications позволяет включить или выключить уведомления.
func (s *ChatUsecaseImpl) SetChatNotofications(ctx context.Context, chatUUID uuid.UUID, userId uuid.UUID, value bool) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	role, err := s.repository.GetUserRoleInChat(ctx, userId, chatUUID)
	if err != nil {
		return err
	}
	if role == NotInChat {
		log.Printf("Пользователь %v не состоит в чате %v", userId, chatUUID)
		return &customerror.NoPermissionError{
			User: userId.String(),
			Area: fmt.Sprintf("Пользователь %v не состоит в чате %v", userId, chatUUID),
		}
	}

	err = s.repository.SetChatNotofications(ctx, chatUUID, userId, value)

	return err
}
