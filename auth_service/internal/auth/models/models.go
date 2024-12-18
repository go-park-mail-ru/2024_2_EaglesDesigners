package models

import (
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey    contextKey = "userId"
	UserKey      contextKey = "user"
	MuxParamsKey contextKey = "muxParams"
)

var (
	ErrUserAlreadyExists = errors.New("a user with that username already exists")
	ErrTokenExpired      = errors.New("token expired")
)

// @Schema.
type AuthReqDTO struct {
	Username string `json:"username" example:"user11" valid:"minstringlength(6),matches(^[a-zA-Z0-9_]+$)"`
	Password string `json:"password"  example:"12345678" valid:"minstringlength(8),matches(^[a-zA-Z0-9_]+$)"`
}

// @Schema.
type RegisterReqDTO struct {
	Username string `json:"username" example:"killer1994" valid:"minstringlength(6),matches(^[a-zA-Z0-9_]+$)"`
	Name     string `json:"name" example:"Vincent Vega" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	Password string `json:"password" example:"go_do_a_crime" valid:"minstringlength(8),matches(^[a-zA-Z0-9_]+$)"`
}

// @Schema.
type RegisterRespDTO struct {
	Message string          `json:"message" example:"Registration successful" valid:"matches(^[a-zA-Zа-яА-Я0-9 ]+$)"`
	User    UserDataRespDTO `json:"user" valid:"-"`
}

// @Schema.
type AuthRespDTO struct {
	User UserDataRespDTO `json:"user" valid:"-"`
}

// @Schema.
type SignupRespDTO struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}

// @Schema.
type UserRespDTO struct {
	ID       uuid.UUID `json:"id" example:"f0364477-bfd4-496d-b639-d825b009d509" valid:"uuid"`
	Username string    `json:"username" example:"mavrodi777" valid:"minstringlength(6),matches(^[a-zA-Z0-9_]+$)"`
	Name     string    `json:"name" example:"Мафиозник" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	Password string    `json:"password" example:"1234567890" valid:"minstringlength(8),matches(^[a-zA-Z0-9_]+$)"`
	Version  int64     `json:"version" example:"1" valid:"int"`
}

// @Schema.
type UserDataRespDTO struct {
	ID        uuid.UUID `json:"id" example:"f0364477-bfd4-496d-b639-d825b009d509" valid:"uuid"`
	Username  string    `json:"username" example:"user12" valid:"minstringlength(6),matches(^[a-zA-Z0-9_]+$)"`
	Name      string    `json:"name" example:"Dr Peper" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	AvatarURL *string   `json:"avatarURL" example:"/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"-"`
}

type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
	Version  int64
}

type UserData struct {
	ID        uuid.UUID
	Username  string
	Name      string
	AvatarURL string
}

type UserDAO struct {
	ID        uuid.UUID
	Username  string
	Name      string
	Password  string
	Version   int64
	AvatarURL *string
}

type CsrfDTO struct {
	Token string `json:"csrf"`
}

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
