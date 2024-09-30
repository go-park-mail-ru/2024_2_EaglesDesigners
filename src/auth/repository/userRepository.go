package repository

import (
	"errors"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

var users = map[string]model.User{
	"user1": {
		ID:       1,
		Username: "user1",
		Name:     "Бал Матье",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user2": {
		ID:       2,
		Username: "user2",
		Name:     "Жабка Пепе",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user3": {
		ID:       3,
		Username: "user3",
		Name:     "Dr Peper",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user4": {
		ID:       4,
		Username: "user4",
		Name:     "Vincent Vega",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
}

func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	user, exists := users[username]
	if !exists {
		return user, errors.New("user does not exist")
	}
	return user, nil
}

func (r *UserRepository) CreateUser(username, name, password string) error {
	if _, exists := users[username]; exists {
		return errors.New("the user already exists")
	}
	users[username] = model.User{ID: int64(len(users)) + 1, Username: username, Name: name, Password: password, Version: 0}
	log.Println("created user:", users[username].ID, users[username].Username, users[username].Name)

	return nil
}
