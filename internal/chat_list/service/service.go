package service

import (
	"errors"

	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
)

type ChatService struct {
	tokenService auth.TokenService
	repository   chatlist.ChatRepository
}

func NewChatService(tokenService auth.TokenService, repository chatlist.ChatRepository) *ChatService {
	return &ChatService{
		tokenService: tokenService,
		repository:   repository,
	}
}

func (s *ChatService) GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error) {

	user, err := s.tokenService.GetUserByJWT(cookie)
	if err != nil {
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	return s.repository.GetUserChats(&user), nil
}
