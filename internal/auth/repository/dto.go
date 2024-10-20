package repository

import "github.com/google/uuid"

// @Schema
type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
	Version  int64
}
