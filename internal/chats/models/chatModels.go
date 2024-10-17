package model

import (
	"path"

	"github.com/google/uuid"
)

// @Schema
type Chat struct {
	ChatId   uuid.UUID
	ChatName string
	// @Enum [personalMessages, group, channel]
	ChatType    string 
	AvatarURL   string 
	ChatURLName string
}


type ChatDTO struct {
	ChatName string `json:"chatName" example:"Чат с пользователем 2"`
	CountOfUsers int `json:"countOfUsers"`
	// @Enum [personalMessages, group, channel]
	ChatType    string `json:"chatType" example:"personalMessages"`
	LastMessage string `json:"lastMessage" example:"Когда за кофе?"`
	AvatarURL   string `json:"avatarURL" example:"https://yandex-images.clstorage.net/bVLC53139/"`
}

type ChatDAO struct {
	ChatId uuid.UUID
	ChatName string
	ChatTypeId int
	AvatarURL string
	ChatURLName string
}


type ChatsDTO struct {
	Chats []ChatDTO `json:"chats"`
}

