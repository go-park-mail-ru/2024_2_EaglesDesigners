package repository

import "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"

type MessagesRepository interface {
	GetMessages(page int, chatId int) ([]models.Message, error)
	AddMessage(models.Message) error
}