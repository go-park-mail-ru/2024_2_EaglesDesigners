package repository

import (
	"errors"
	"log"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

var users = map[string]User{
	"user11": {
		ID:       1,
		Username: "user11",
		Name:     "Бал Матье",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user22": {
		ID:       2,
		Username: "user22",
		Name:     "Жабка Пепе",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user33": {
		ID:       3,
		Username: "user33",
		Name:     "Dr Peper",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user44": {
		ID:       4,
		Username: "user44",
		Name:     "Vincent Vega",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
}

func (r *Repository) GetUserByUsername(username string) (User, error) {
	user, exists := users[username]
	if !exists {
		log.Println("Пользователь не найден в базе данных")
		return user, errors.New("user does not exist")
	}
	log.Printf("Пользователь c id %d найден", user.ID)
	return user, nil
}

func (r *Repository) CreateUser(username, name, password string) error {
	if _, exists := users[username]; exists {
		return errors.New("the user already exists")
	}
	users[username] = User{ID: int64(len(users)) + 1, Username: username, Name: name, Password: password, Version: 0}
	log.Println("created user:", users[username].ID, users[username].Username, users[username].Name)

	return nil
}
