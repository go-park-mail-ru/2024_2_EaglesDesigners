package login

import (
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/controller"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/service"
)

func SetupController() *controller.LoginController {
	userRepo := &repository.UserRepository{}
	userService := &service.UserService{
		ILoginRepository: userRepo,
	}
	loginController := &controller.LoginController{
		ILoginService: userService,
	}
	return loginController
}
