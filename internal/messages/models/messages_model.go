package models

import "time"

type Message struct {
	MessageId int    `json:"messageId" example:"1"`
	AuthorID  int    `json:"authorID" exameple:"2"`
	Message   string `json:"text" example:"тут много текста"`
	SentAt time.Time `json:"datetime" example:"2024-04-13T08:30:00Z"`
}

type MessagesArrayDTO struct {
	Messages []Message
}
