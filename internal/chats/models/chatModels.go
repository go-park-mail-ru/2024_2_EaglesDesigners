package model

import (
	"encoding/json"

	"github.com/google/uuid"

	messageModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
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

type ChatDTOOutput struct {
	ChatId       uuid.UUID `json:"chatId"`
	ChatName     string    `json:"chatName" example:"Чат с пользователем 2"`
	CountOfUsers int       `json:"countOfUsers"`
	// @Enum [personalMessages, group, channel]
	ChatType    string               `json:"chatType" example:"personalMessages"`
	LastMessage messageModel.Message `json:"lastMessage"`
	// фото в формате base64
	AvatarBase64 string `json:"avatarBase64"`
}

type ChatDTOInput struct {
	ChatName     string      `json:"chatName" example:"Чат с пользователем 2"`
	ChatType     string      `json:"chatType" example:"personalMessages"`
	AvatarBase64 string      `json:"avatarBase64"`
	UsersToAdd   []uuid.UUID `json:"usersToAdd"`
}

func (chat ChatDTOOutput) MarshalBinary() ([]byte, error) {
	return json.Marshal(chat)
}

func (chat *ChatDTOOutput) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, chat)
}

type ChatDAO struct {
	ChatId      uuid.UUID
	ChatName    string
	ChatTypeId  int
	AvatarURL   string
	ChatURLName string
}

type ChatsDTO struct {
	Chats []ChatDTOOutput `json:"chats"`
}

func СhatToChatDTO(chat Chat, countOfUsers int, lastMessage messageModel.Message, AvatarBase64 string) ChatDTOOutput {
	return ChatDTOOutput{
		ChatId:       chat.ChatId,
		ChatName:     chat.ChatName,
		CountOfUsers: countOfUsers,
		ChatType:     chat.ChatType,
		LastMessage:  lastMessage,
		AvatarBase64: AvatarBase64,
	}
}

type AddUsersIntoChatDTO struct {
	ChatId  uuid.UUID   `json:"chatId"`
	UsersId []uuid.UUID `json:"usersId"`
}
