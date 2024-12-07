package models

import (
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// @Schema
type Message struct {
	MessageId    uuid.UUID               `json:"messageId" example:"1" valid:"-"`
	AuthorID     uuid.UUID               `json:"authorID" exameple:"2" valid:"-"`
	BranchID     *uuid.UUID              `json:"branchId" exameple:"2" valid:"-"`
	Message      string                  `json:"text" example:"тут много текста" valid:"-"`
	SentAt       time.Time               `json:"datetime" example:"2024-04-13T08:30:00Z" valid:"-"`
	ChatId       uuid.UUID               `json:"chatId" valid:"-"`
	IsRedacted   bool                    `json:"isRedacted" valid:"-"`
	MessageType  string                  `json:"message_type" valid:"-" example:"informational"`
	Files        []multipart.File        `json:"-" valid:"-"`
	FilesHeaders []*multipart.FileHeader `json:"-" valid:"-"`
	FilesURLs    []string                `json:"files" valid:"-" example:"[url1, url2, url3]"`
	Photos        []multipart.File        `json:"-" valid:"-"`
	PhotosHeaders []*multipart.FileHeader `json:"-" valid:"-"`
	PhotosURLs    []string                `json:"photos" valid:"-" example:"[url1, url2, url3]"`
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

type WebScoketDTO struct {
	MsgType MsgType     `json:"messageType"`
	Payload interface{} `json:"payload"`
}
