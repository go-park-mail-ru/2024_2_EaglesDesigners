package usecase

import (
	"context"
	"errors"
	"log"

	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/base64helper"
	"github.com/google/uuid"
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
	tokenUsecase *usecase.Usecase
	repository   chatlist.ChatRepository
}

func NewChatUsecase(tokenService *usecase.Usecase, repository chatlist.ChatRepository) ChatUsecase {
	return &ChatUsecaseImpl{
		tokenUsecase: tokenService,
		repository:   repository,
	}
}

func (s *ChatUsecaseImpl) GetChats(ctx context.Context, cookie []*http.Cookie, pageNum int) ([]chatModel.Chat, error) {

	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	chats, err := s.repository.GetUserChats(user.ID, pageNum)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *ChatUsecaseImpl) AddUsersIntoChat(ctx context.Context, cookie []*http.Cookie, user_ids []uuid.UUID, chat_id uuid.UUID) error {
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	role, err := s.repository.GetUserRoleInChat(user.ID, chat_id)
	if err != nil {
		return err
	}
	switch role {
	case admin, owner:
		log.Printf("Начато добавление пользователей в чат %v пользователем %v", chat_id, user.ID)
		for _, id := range user_ids {
			s.repository.AddUserIntoChat(id, chat_id, none)
		}
		log.Printf("Участники добавлены в чат %v пользователем %v", chat_id, user.ID)
		return nil
	}

	return errors.New("Участники не добавлены")
}

func (s *ChatUsecaseImpl) AddNewChat(ctx context.Context, cookie []*http.Cookie, chat chatModel.ChatDTO) error {
	user, err := s.tokenUsecase.GetUserByJWT(ctx, cookie)
	if err != nil {
		return errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	photoPath, err := base64helper.SavePhotoBase64(chat.AvatarBase64)

	if err != nil {
		log.Printf("Не удалось сохранить фото^ %v", err)
		return err
	}

	chatId := uuid.New()

	newChat := chatModel.Chat{
		ChatId:    chatId,
		ChatName:  chat.ChatName,
		ChatType:  chat.ChatType,
		AvatarURL: photoPath.String(),
	}

	if chat.ChatType == channel {
		newChat.ChatURLName = chat.ChatName + "_" + uuid.NewString()
	}

	// создание чата
	err = s.repository.CreateNewChat(newChat)
	if err != nil {
		log.Printf("Не удалось сохнанить чат: %v", err)
		return err
	}

	// добавление владельца
	err = s.repository.AddUserIntoChat(user.ID, chatId, owner)

	if err != nil {
		log.Printf("Не удалось добавить пользователя в чат: %v", err)
		return err
	}

	return nil
}
