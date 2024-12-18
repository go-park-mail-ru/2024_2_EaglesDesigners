package events

import (
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

const (
	UpdateChat          = "updateChat"
	DeleteChat          = "deleteChat"
	NewChat             = "newChat"
	DeleteUsersFromChat = "delUsers"
	AddNewUsersInChat   = "addUsers"

	// пользователь стал онлайн.
	AddWebcosketUser = "addWebSocketUser"
)

const (
	DeleteMessage = "deleteMessage"
	NewMessage    = "newMessage"
	UpdateMessage = "updateMessage"
)

type MessageEvent struct {
	Action  string  `json:"action"`
	Message Message `json:"payload"`
}
type Message struct {
	MessageId     uuid.UUID               `json:"messageId" example:"1" valid:"-"`
	AuthorID      uuid.UUID               `json:"authorID" exameple:"2" valid:"-"`
	BranchID      *uuid.UUID              `json:"branchId" exameple:"2" valid:"-"`
	Message       string                  `json:"text" example:"тут много текста" valid:"-"`
	SentAt        time.Time               `json:"datetime" example:"2024-04-13T08:30:00Z" valid:"-"`
	ChatId        uuid.UUID               `json:"chatId" valid:"-"`
	IsRedacted    bool                    `json:"isRedacted" valid:"-"`
	ChatIdParent  uuid.UUID               `json:"parent_chat_id" valid:"-"`
	Files         []multipart.File        `json:"-" valid:"-"`
	FilesHeaders  []*multipart.FileHeader `json:"-" valid:"-"`
	FilesDTO      []Payload               `json:"files" valid:"-"`
	Photos        []multipart.File        `json:"-" valid:"-"`
	PhotosHeaders []*multipart.FileHeader `json:"-" valid:"-"`
	PhotosDTO     []Payload               `json:"photos" valid:"-"`
	Sticker       string                  `json:"sticker" valid:"-" example:"/files/675f2ea013dbaf51a93aa2d3"`
}
type Payload struct {
	URL      string `json:"url" example:"url" valid:"-"`
	Filename string `json:"filename" example:"image.png" valid:"-"`
	Size     int64  `json:"size" example:"10500" valid:"-"`
}

func SerializeMessageEvent(event MessageEvent) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeMessageEvent(data []byte) (MessageEvent, error) {
	var event MessageEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return MessageEvent{}, err
	}
	return event, nil
}

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
