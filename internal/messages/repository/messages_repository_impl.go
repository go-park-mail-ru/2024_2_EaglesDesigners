package repository

import "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"

type MessagesRepositoryImpl struct {
}

func (r *MessagesRepositoryImpl) GetMessages(page int, chatId int) ([]models.Message, error) {
	return nil, nil
}

func (r *MessagesRepositoryImpl) AddMessage(models.Message) error {
	return nil
}
