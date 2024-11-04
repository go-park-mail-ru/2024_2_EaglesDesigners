package model

import (
	"encoding/json"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
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

// @Schema
type ChatDTOOutput struct {
	ChatId       uuid.UUID      `json:"chatId"`
	ChatName     string         `json:"chatName" example:"Чат с пользователем 2"`
	CountOfUsers int            `json:"countOfUsers" example:"52"`
	ChatType     string         `json:"chatType" example:"personal"`
	LastMessage  models.Message `json:"lastMessage"`
	AvatarPath   string         `json:"avatarPath"`
}

type ChatDTOInput struct {
	ChatName   string          `json:"chatName" example:"Чат с пользователем 2"`
	ChatType   string          `json:"chatType" example:"personalMessages"`
	UsersToAdd []uuid.UUID     `json:"usersToAdd" example:"uuid1,uuid2"`
	Avatar     *multipart.File `json:"-"`
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

type ChatUpdate struct {
	ChatName string          `json:"chatName" example:"Чат с пользователем 2"`
	Avatar   *multipart.File `json:"-"`
}

func СhatToChatDTO(chat Chat, countOfUsers int, lastMessage models.Message) ChatDTOOutput {
	return ChatDTOOutput{
		ChatId:       chat.ChatId,
		ChatName:     chat.ChatName,
		CountOfUsers: countOfUsers,
		ChatType:     chat.ChatType,
		LastMessage:  lastMessage,
		AvatarPath:   chat.AvatarURL,
	}
}

type AddUsersIntoChatDTO struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2"`
}

type AddedUsersIntoChatDTO struct {
	AddedUsers []uuid.UUID `json:"addedUser" example:"uuid1,uuid2"`
}

type DeleteUsersFromChatDTO struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2"`
}

type DeletdeUsersFromChatDTO struct {
	DeletedUsers []uuid.UUID `json:"deletedUsers" example:"uuid1,uuid2"`
}
