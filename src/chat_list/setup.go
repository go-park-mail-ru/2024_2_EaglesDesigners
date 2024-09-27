package chatlist

import (
	userRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/repository"
	authService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
	chatController "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/controller"
	chatRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/repository"
	chatService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/service"
)

func SetupController() *chatController.ChatController {
	userRepository := userRepo.NewUserRepository()
	chatRepository := chatRepo.NewChatRepository()
	tokenService := authService.NewTokenService(*userRepository)
	chatService := chatService.NewChatService(*tokenService, *chatRepository)
	controller := chatController.NewChatController(*chatService)

	return controller
}
