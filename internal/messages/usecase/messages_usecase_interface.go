package usecase

import "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"

type MessagesUsecase interface {
	SendMessage(userId int, chatId int) error
	GetMessages(page int, chatId int) ([]models.MessagesArrayDTO, error)
}