package model

import (
	"encoding/json"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"
)

// @Schema.
type Chat struct {
	ChatId   uuid.UUID
	ChatName string
	// @Enum [personalMessages, group, channel]
	ChatType string
	// путь до фото в папке Messages
	AvatarURL string
	// типо неймтаг канала
	ChatURLName       string
	SendNotifications bool
}

// @Schema
//
//easyjson:json
type ChatDTOOutput struct {
	ChatId            uuid.UUID      `json:"chatId" example:"08a0f350-e122-467b-8ba8-524d2478b56e" valid:"-"`
	ChatName          string         `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	CountOfUsers      int            `json:"countOfUsers" example:"52" valid:"int"`
	ChatType          string         `json:"chatType" example:"personal" valid:"in(personal|group|channel)"`
	LastMessage       models.Message `json:"lastMessage" valid:"-"`
	AvatarPath        string         `json:"avatarPath"  example:"/uploads/chat/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"-"`
	SendNotifications bool           `json:"send_notifications" valid:"-"`
}

// для сортировки возвращаемого списка по убыванию.
type ByLastMessage []ChatDTOOutput

func (a ByLastMessage) Len() int { return len(a) }

func (a ByLastMessage) Less(i, j int) bool {
	return a[i].LastMessage.SentAt.After(a[j].LastMessage.SentAt)
}

func (a ByLastMessage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//easyjson:json
type ChatDTOInput struct {
	ChatName     string                `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	ChatType     string                `json:"chatType" example:"personalMessages" valid:"in(personal|group|channel)"`
	UsersToAdd   []uuid.UUID           `json:"usersToAdd" example:"uuid1,uuid2" valid:"-"`
	Avatar       *multipart.File       `json:"-" valid:"-"`
	AvatarHeader *multipart.FileHeader `json:"-" valid:"-"`
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

//easyjson:json
type ChatsDTO struct {
	Chats []ChatDTOOutput `json:"chats" valid:"-"`
}

//easyjson:json
type ChatUpdate struct {
	ChatName     string                `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	Avatar       *multipart.File       `json:"-" valid:"-"`
	AvatarHeader *multipart.FileHeader `json:"-" valid:"-"`
}

//easyjson:json
type ChatUpdateOutput struct {
	ChatName string `json:"chatName" example:"Чат с пользователем 2" valid:"-"`
	Avatar   string `json:"updatedAvatarPath" example:"/uploads/chat/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"-"`
}

func ChatToChatDTO(chat Chat, countOfUsers int, lastMessage models.Message) ChatDTOOutput {
	return ChatDTOOutput{
		ChatId:            chat.ChatId,
		ChatName:          chat.ChatName,
		CountOfUsers:      countOfUsers,
		ChatType:          chat.ChatType,
		LastMessage:       lastMessage,
		AvatarPath:        chat.AvatarURL,
		SendNotifications: chat.SendNotifications,
	}
}

//easyjson:json
type AddUsersIntoChatDTO struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2" valid:"-"`
}

//easyjson:json
type AddedUsersIntoChatDTO struct {
	AddedUsers    []uuid.UUID `json:"addedUser" example:"uuid1,uuid2" valid:"-"`
	NotAddedUsers []uuid.UUID `json:"notAddedUser" example:"uuid1,uuid2" valid:"-"`
}

//easyjson:json
type DeleteUsersFromChatDTO struct {
	UsersId []uuid.UUID `json:"usersId" example:"uuid1,uuid2" valid:"-"`
}

//easyjson:json
type DeletdeUsersFromChatDTO struct {
	DeletedUsers []uuid.UUID `json:"deletedUsers" example:"uuid1,uuid2" valid:"-"`
}

//easyjson:json
type ChatInfoDTO struct {
	Role              string           `json:"role" example:"owner" valid:"in(admin|owner|none)"`
	Users             []UserInChatDTO  `json:"users" valid:"-"`
	Messages          []models.Message `json:"messages" valid:"-"`
	SendNotifications bool             `json:"send_notifications" valid:"-"`
	Files             []models.Payload `json:"files" valid:"-"`
	Photos            []models.Payload `json:"photos" valid:"-"`
}

//easyjson:json
type UserInChatDTO struct {
	ID         uuid.UUID `json:"id" example:"f0364477-bfd4-496d-b639-d825b009d509" valid:"uuid"`
	Username   string    `json:"username" example:"mavrodi777" valid:"minstringlength(6),matches(^[a-zA-Z0-9_]+$)"`
	Name       *string   `json:"name" example:"Vincent Vega" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$),optional"`
	AvatarPath *string   `json:"avatarURL" example:"/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"-"`
	Role       *string   `json:"role" example:"owner" valid:"in(admin|owner|none),optional"`
}

type UserInChatDAO struct {
	ID         uuid.UUID
	Username   string
	Name       *string
	AvatarPath *string
	Role       *int
}

//easyjson:json
type Event struct {
	Action string      `json:"action"`
	ChatId uuid.UUID   `json:"chatId"`
	Users  []uuid.UUID `json:"users"`
}

func SerializeEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeEvent(data []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

//easyjson:json
type AddBranch struct {
	ID uuid.UUID `json:"id" example:"f0364477-bfd4-496d-b639-d825b009d509" valid:"uuid"`
}

//easyjson:json
type SearchChatsDTO struct {
	UserChats      []ChatDTOOutput `json:"user_chats" valid:"-"`
	GlobalChannels []ChatDTOOutput `json:"global_channels" valid:"-"`
}
