package service

import (
	"../repository"
	"errors"
	"fmt"
	"net/http"
)

func isAuthorized(cookie []*http.Cookie) bool {

	return true
}
func getUserByJWT(cookie []*http.Cookie) repository.User {
	return repository.User{UserId: 1, Username: "Олег", Password: "1123"}
}

func GetChats(cookie []*http.Cookie) ([]repository.Chat, error) {
	fmt.Println("yes")

	if isAuthorized(cookie) {
		user := getUserByJWT(cookie)
		chats := repository.GetUserChats(&user)
		fmt.Println(chats)

		return chats, nil
	} else {
		return nil, errors.New("Пользователь не авторизован")
	}

}
