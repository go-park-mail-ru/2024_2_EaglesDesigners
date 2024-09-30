package service

import (
	"errors"
	"log"

	"net/http"

	userService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/repository"
)

type ChatService struct {
	tokenService userService.TokenService
	repository   repository.ChatRepository
}

func NewChatService(tokenService userService.TokenService, repository repository.ChatRepository) *ChatService {
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
