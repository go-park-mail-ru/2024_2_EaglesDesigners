package model

import (
	"github.com/google/uuid"
)

// @Schema
type Chat struct {
	ChatId   uuid.UUID
	ChatName string
	// @Enum [personalMessages, group, channel]
	ChatType string
	// путь до фото в папке Messages
	AvatarURL string
	// типо неймтаг канала
	ChatURLName string
}

type ChatDTO struct {
	ChatName     string `json:"chatName" example:"Чат с пользователем 2"`
	CountOfUsers int    `json:"countOfUsers"`
	// @Enum [personalMessages, group, channel]
	ChatType    string `json:"chatType" example:"personalMessages"`
	LastMessage string `json:"lastMessage" example:"Когда за кофе?"`
	// фото в формате base64
	AvatarBase64 string `json:"avatarBase64"`
}

type ChatDAO struct {
	ChatId      uuid.UUID
	ChatName    string
	ChatTypeId  int
	AvatarURL   string
	ChatURLName string
}

type ChatsDTO struct {
	Chats []ChatDTO `json:"chats"`
}

func СhatToChatDTO(chat Chat, countOfUsers int, lastMessage string, AvatarBase64 string) ChatDTO {
	return ChatDTO{
		ChatName:     chat.ChatName,
		CountOfUsers: countOfUsers,
		ChatType:     chat.ChatType,
		LastMessage:  lastMessage,
		AvatarBase64: AvatarBase64,
	}
}
