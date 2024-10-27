package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId  uuid.UUID `json:"messageId" example:"1"`
	AuthorID   uuid.UUID `json:"authorID" exameple:"2"`
	AuthorName string    `json:"authorName"`
	Message    string    `json:"text" example:"тут много текста"`
	SentAt     time.Time `json:"datetime" example:"2024-04-13T08:30:00Z"`
	IsRedacted bool      `json:"isRedacted"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// Custom unmarshaling for Message
func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type MessagesArrayDTO struct {
	Messages []Message `json:"messages"`
}

type MessagesArrayDTOOutput struct {
	IsNew    bool      `json:"isNew"`
	Messages []Message `json:"messages"`
}

type MessageDTOInput struct {
	Disconnect bool   `json:"disconnect"`
	Message    string `json:"message"`
}
