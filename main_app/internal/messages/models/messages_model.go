package models

import (
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// @Schema
//
//easyjson:json
type Message struct {
	MessageId     uuid.UUID               `json:"messageId" example:"1" valid:"-"`
	AuthorID      uuid.UUID               `json:"authorID" exameple:"2" valid:"-"`
	BranchID      *uuid.UUID              `json:"branchId" exameple:"2" valid:"-"`
	Message       string                  `json:"text" example:"тут много текста" valid:"-"`
	SentAt        time.Time               `json:"datetime" example:"2024-04-13T08:30:00Z" valid:"-"`
	ChatId        uuid.UUID               `json:"chatId" valid:"-"`
	IsRedacted    bool                    `json:"isRedacted" valid:"-"`
	MessageType   string                  `json:"message_type" valid:"-" example:"with_payload"`
	ChatIdParent  uuid.UUID               `json:"parent_chat_id" valid:"-"`
	Files         []multipart.File        `json:"-" valid:"-"`
	FilesHeaders  []*multipart.FileHeader `json:"-" valid:"-"`
	FilesDTO      []Payload               `json:"files" valid:"-"`
	Photos        []multipart.File        `json:"-" valid:"-"`
	PhotosHeaders []*multipart.FileHeader `json:"-" valid:"-"`
	PhotosDTO     []Payload               `json:"photos" valid:"-"`
	Sticker       string                  `json:"sticker" valid:"-" example:"/files/675f2ea013dbaf51a93aa2d3"`
}

//easyjson:json
type Payload struct {
	URL      string `json:"url" example:"url" valid:"-"`
	Filename string `json:"filename" example:"image.png" valid:"-"`
	Size     int64  `json:"size" example:"10500" valid:"-"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// Custom unmarshaling for Message
func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type MessageInput struct {
	Message string `json:"text" example:"тут много текста" valid:"-"`
}

//easyjson:json
type MessagesArrayDTO struct {
	Messages []Message `json:"messages" valid:"-"`
}

type MessagesArrayDTOOutput struct {
	IsNew    bool      `json:"isNew" valid:"-"`
	Messages []Message `json:"messages" valid:"-"`
}

type MessageDTOInput struct {
	Disconnect bool   `json:"disconnect" valid:"-"`
	Message    string `json:"message" valid:"-"`
}

type MsgType string

const (
	NewMessage      MsgType = "message"
	FeatUserInChat  MsgType = "featUserInChat"
	DelUserFromChat MsgType = "delUserFromChat"
)

//easyjson:skip
type WebScoketDTO struct {
	MsgType MsgType     `json:"messageType"`
	Payload interface{} `json:"payload"`
}
