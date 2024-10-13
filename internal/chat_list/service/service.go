package service

import (
	"errors"

	"net/http"

	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type ChatService struct {
	tokenUsecase usecase.Usecase
	repository   chatlist.ChatRepository
}

func NewChatService(tokenService usecase.Usecase, repository chatlist.ChatRepository) *ChatService {
	return &ChatService{
		tokenUsecase: tokenService,
		repository:   repository,
	}
}

func (s *ChatService) GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error) {

	user, err := s.tokenUsecase.GetUserByJWT(cookie)
	if err != nil {
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	return s.repository.GetUserChats(&user), nil
}
