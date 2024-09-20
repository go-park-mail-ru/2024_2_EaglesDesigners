package service

import (
	interfaces "github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/interface"
)

type UserService struct {
	interfaces.ILoginRepository
}

func (service *UserService) Authenticate(username, password string) bool {
	user := service.FindByUsername(username)
	return user.Password == password
}
