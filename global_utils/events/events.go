package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageEvent struct {
	Action  string               `json:"action"`
	Message Message `json:"payload"`
}
type Message struct {
	MessageId  uuid.UUID  `json:"messageId" example:"1" valid:"-"`
	AuthorID   uuid.UUID  `json:"authorID" exameple:"2" valid:"-"`
	BranchID   *uuid.UUID `json:"branchId" exameple:"2" valid:"-"`
	Message    string     `json:"text" example:"тут много текста" valid:"-"`
	SentAt     time.Time  `json:"datetime" example:"2024-04-13T08:30:00Z" valid:"-"`
	ChatId     uuid.UUID  `json:"chatId" valid:"-"`
	IsRedacted bool       `json:"isRedacted" valid:"-"`
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