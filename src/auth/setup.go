package auth

import (
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/controller"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
)

func SetupController() *controller.AuthController {
	userRepo := repository.NewUserRepository()
	tokenService := service.NewTokenService(*userRepo)
	authService := service.NewAuthService(*userRepo, *tokenService)
	authController := controller.NewAuthController(*authService, *tokenService)
	return authController
}
