package models

import "github.com/google/uuid"

type Question struct {
	QuestionId   uuid.UUID
	QuestionText string
	QuestionType string
}

type Answer struct {
	AnswerId   uuid.UUID
	QuestionId   uuid.UUID
	UserId uuid.UUID
	TextAnswer string
	NumericAnswer int
}
