package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/google/uuid"
)

const (
	admin = "admin"
	none  = "none"
	owner = "owner"
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
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	photoBytes, err := base64.RawStdEncoding.DecodeString(chat.AvatarBase64)

	return nil
}

func savePhotoBase64(base64Photo string) (uuid.UUID, error) {
	photoBytes, err := base64.RawStdEncoding.DecodeString(base64Photo)
	if err != nil {
		log.Printf("Ну удалось расшифровать фото: %v \n", err)
		return uuid.Nil, err
	}

	permissions := 777 // or whatever you need
	filenameUUID := uuid.New()

	err = os.WriteFile("../../images/" + filenameUUID.String(), photoBytes, permissions)

	if err != nil {
		log.Printf("Unable to write into file %v: %v", filenameUUID, err)
		return uuid.Nil, err
	}
	return filenameUUID, nil
}
