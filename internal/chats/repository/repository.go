package repository

import (
	"log"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)


type ChatRepositoryImpl struct {
	
}

func NewChatRepository() ChatRepository {
	return &ChatRepositoryImpl{}
}

func (r *ChatRepositoryImpl) GetUserChats(user *userModel.User) []chatModel.Chat {
	log.Printf("Поиск чатов пользователя %d", user.ID)

	chats, ok := keys[user.ID]

	if !ok {
		log.Printf("Чаты пользователья %d не найдены", user.ID)
		return []chatModel.Chat{}
	}
	log.Printf("Найдено чатов пользователя %d: %d ", user.ID, len(chats))
	return chats
}

func (r* ChatRepositoryImpl) IsUserInChat(userId int, chatId int) bool {
	// идем в бд по двум полям: если есть то тру

	//а пока так:
	chats, ok := keys[int64(userId)] 
	if !ok {
		log.Printf("Чаты пользователья %d не найдены", userId)
		return false
	}
	for _, chat := range chats {
		if chat.ChatId == chatId {
			return true
		}
	}
	return false
}