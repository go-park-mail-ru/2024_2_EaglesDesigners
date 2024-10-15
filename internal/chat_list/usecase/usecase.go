package usecase

import (
	"errors"

	"net/http"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
	chatlist "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/repository"
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
