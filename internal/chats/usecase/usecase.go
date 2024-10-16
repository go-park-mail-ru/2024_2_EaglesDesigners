package usecase

import (
	"errors"
	"log"

	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
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

func (s *ChatUsecaseImpl) GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error) {

	user, err := s.tokenUsecase.GetUserByJWT(cookie)
	if err != nil {
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}
	
	return s.repository.GetUserChats(&user), nil
}

func (s *ChatUsecaseImpl) CanUserWriteInChat(userId int, chatId int) bool {
	// проверяем состоит ли пользователь в чате
	ok := s.repository.IsUserInChat(userId, chatId)
	if !ok {
		log.Printf("Пользователь %d не состоит в чате %d", userId, chatId)
		return false
	}

	// проверяем есть ли у пользователя права писать в чат
	// role := s.repository.GetUserRole(userId, chatId)
	// if !role return false

	return true
}

