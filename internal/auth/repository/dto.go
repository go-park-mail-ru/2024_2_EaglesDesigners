package repository

import "github.com/google/uuid"

// @Schema
type User struct {
	ID       uuid.UUID `json:"id" example:"1"`
	Username string    `json:"username" example:"mavrodi777"`
	Name     string    `json:"name" example:"Мафиозник"`
	Password string    `json:"password" example:"1234567890"`
	Version  int64     `json:"version" example:"1"`
}
