package models

import "github.com/google/uuid"

type contextKey string

const (
	UserIDKey    contextKey = "userId"
	UserKey      contextKey = "user"
	MuxParamsKey contextKey = "muxParams"
)

// @Schema
type AuthReqDTO struct {
	Username string `json:"username" example:"user11"`
	Password string `json:"password"  example:"12345678"`
}

// @Schema
type RegisterReqDTO struct {
	Username string `json:"username" example:"killer1994"`
	Name     string `json:"name" example:"Vincent Vega"`
	Password string `json:"password" example:"go_do_a_crime"`
}

// @Schema
type RegisterRespDTO struct {
	Message string          `json:"message" example:"Registration successful"`
	User    UserDataRespDTO `json:"user"`
}

// @Schema
type AuthRespDTO struct {
	User UserDataRespDTO `json:"user"`
}

// @Schema
type SignupRespDTO struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}

// @Schema
type UserRespDTO struct {
	ID       uuid.UUID `json:"id" example:"1"`
	Username string    `json:"username" example:"mavrodi777"`
	Name     string    `json:"name" example:"Мафиозник"`
	Password string    `json:"password" example:"1234567890"`
	Version  int64     `json:"version" example:"1"`
}

// @Schema
type UserDataRespDTO struct {
	ID        uuid.UUID `json:"id" example:"2"`
	Username  string    `json:"username" example:"user12"`
	Name      string    `json:"name" example:"Dr Peper"`
	AvatarURL *string   `json:"avatarURL" example:"/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png"`
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
	AvatarURL *string
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
