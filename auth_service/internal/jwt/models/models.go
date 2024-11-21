package models

import "github.com/google/uuid"

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub     string    `json:"sub"`
	Name    string    `json:"name"`
	ID      uuid.UUID `json:"id"`
	Version int64     `json:"vrs"`
	Exp     int64     `json:"exp"`
}

type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
	Version  int64
}

type UserData struct {
	ID       uuid.UUID `json:"id" example:"2"`
	Username string    `json:"username" example:"user12"`
	Name     string    `json:"name" example:"Dr Peper"`
}
