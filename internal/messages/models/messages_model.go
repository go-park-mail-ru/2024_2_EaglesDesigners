package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// @Schema
type Message struct {
	MessageId  uuid.UUID `json:"messageId" example:"1" valid:"uuid"`
	AuthorID   uuid.UUID `json:"authorID" exameple:"2" valid:"uuid"`
	AuthorName string    `json:"authorName" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	Message    string    `json:"text" example:"тут много текста" valid:"-"`
	SentAt     time.Time `json:"datetime" example:"2024-04-13T08:30:00Z" valid:"-"`
	ChatId     uuid.UUID `json:"chatId" valid:"uuid"`
	IsRedacted bool      `json:"isRedacted" valid:"-"`
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
