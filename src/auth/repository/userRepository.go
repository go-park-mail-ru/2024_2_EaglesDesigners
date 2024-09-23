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
		Password: "5daafb93f7b8c5b1041a5611bb8ca082cd62ea1d0462717a4f2b06e160a531263ddd4ae3f241c6fd0551b97dbad837fb8892db6d9ebba9e98ddf4f6e83f47f3d",
		Salt:     []byte{1, 25, 166, 86, 76, 151, 89, 178, 211, 185, 117, 124, 253, 7, 197, 248},
	},
	"user2": {
		ID:       1,
		Username: "user2",
		Password: "65c07e80473a2c915b1b1fd769610cf4d2c5dd20d07716bf7622d340cba1b9a16e9055a8ff8d2c2a0dfb15212a7506862e04a41075e655f2687048cbdb11592d",
		Salt:     []byte{118, 45, 116, 226, 167, 225, 122, 104, 27, 166, 92, 198, 92, 112, 55, 128},
	},
}

func (repository *UserRepository) GetUserByUsername(username string) model.User {
	return users[username]
	// user, exists := users[username]
	// if exists {
	// 	return user
	// }

	// return model.User{}
}

func (repository *UserRepository) CreateUser(username, password string, salt []byte) error {
	if _, exists := users[username]; exists {
		return errors.New("the user already exists")
	}
	users[username] = model.User{ID: int64(len(users)), Username: username, Password: password, Salt: salt}

	return nil
}
