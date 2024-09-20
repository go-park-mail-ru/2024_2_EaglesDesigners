package repository

import (
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/model"
)

type UserRepository struct{}

var users = []model.User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
	{Username: "user3", Password: "pass3"},
}

func (repository *UserRepository) FindByUsername(username string) model.User {
	for _, user := range users {
		if user.Username == username {
			return user
		}
	}
	return model.User{}
}
