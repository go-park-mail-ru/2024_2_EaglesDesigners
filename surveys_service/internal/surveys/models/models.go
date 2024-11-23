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

type Stat struct {
	QuestionId uuid.UUID
	QuestionType string
	Avarage_numeric string
	AllTextAnswers []string
}

type Servey struct {
	Topic string
	Id string
}