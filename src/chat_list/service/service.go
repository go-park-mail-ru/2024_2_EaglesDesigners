package service

import (
	"errors"
	"fmt"
	"net/http"

	userRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/repository"
	userService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/repository"
)

var tokenService *userService.TokenService

func GetChats(cookie []*http.Cookie) ([]chatModel.Chat, error) {
	fmt.Println("yes")

	user, err := tokenService.GetUserByJWT(cookie)
	if err != nil {
		return []chatModel.Chat{}, errors.New("НЕ УДАЛОСЬ ПОЛУЧИТЬ ПОЛЬЗОВАТЕЛЯ")
	}

	return repository.GetUserChats(&user), nil

}

func init() {
	tokenService = userService.NewTokenService(*userRepository.NewUserRepository())
}
