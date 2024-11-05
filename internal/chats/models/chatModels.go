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
	ChatId       uuid.UUID      `json:"chatId" example:"08a0f350-e122-467b-8ba8-524d2478b56e" valid:"uuid"`
	ChatName     string         `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	CountOfUsers int            `json:"countOfUsers" example:"52" valid:"int"`
	ChatType     string         `json:"chatType" example:"personal" valid:"in(personal|group|channel)"`
	LastMessage  models.Message `json:"lastMessage" valid:",optional"`
	AvatarPath   string         `json:"avatarPath"  example:"/uploads/chat/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"matches(^/uploads/chat/[a-zA-Z0-9\\-]+\\.png$),optional"`
}

type ChatDTOInput struct {
	ChatName   string          `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	ChatType   string          `json:"chatType" example:"personalMessages" valid:"in(personal|group|channel)"`
	UsersToAdd []uuid.UUID     `json:"usersToAdd" example:"uuid1,uuid2" valid:"uuid"`
	Avatar     *multipart.File `json:"-" valid:"-"`
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
	Chats []ChatDTOOutput `json:"chats" valid:"-"`
}

type ChatUpdate struct {
	ChatName string          `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	Avatar   *multipart.File `json:"-"`
}

type ChatUpdateOutput struct {
	ChatName string `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	Avatar   string `json:"updatedAvatarPath" example:"/uploads/chat/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"matches(^/uploads/chat/[a-zA-Z0-9\\-]+\\.png$),optional"`
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
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2" valid:"-"`
}

type AddedUsersIntoChatDTO struct {
	AddedUsers    []uuid.UUID `json:"addedUser" example:"uuid1,uuid2" valid:"-"`
	NotAddedUsers []uuid.UUID `json:"notAddedUser" example:"uuid1,uuid2" valid:"-"`
}

type DeleteUsersFromChatDTO struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2" valid:"-"`
}

type DeletdeUsersFromChatDTO struct {
	DeletedUsers []uuid.UUID `json:"deletedUsers" example:"uuid1,uuid2" valid:"-"`
}

type UsersInChat struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2" valid:"-"`
}
