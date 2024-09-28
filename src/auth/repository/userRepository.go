package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

var users = map[string]model.User{
	"user1": {
		ID:       0,
		Username: "user1",
		Name:     "Бал Матье",
		Password: "edfe8b0b61dd1b86b855e138936db4eb2793df3aaae52c1fb8c5fc08dae72b36f4374a5569e9abf90d85bddf2c6a576b47b5b4ca948e98409c0ca9b3047c9509",
	},
	"user2": {
		ID:       1,
		Username: "user2",
		Name:     "Жабка Пепе",
		Password: "5126c5568550f9ced38b64ce28db77796bfb46927b8c13c2bf83c7b09b539cb9469e0c7c33afbdfdc12fdcaa80bd11c024118cabda0beeafb29c9352e6f017db",
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
	users[username] = model.User{ID: int64(len(users)), Username: username, Name: name, Password: password, Version: 0}

	return nil
}
