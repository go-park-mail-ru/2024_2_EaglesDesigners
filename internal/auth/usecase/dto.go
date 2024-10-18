package usecase

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
	Version  int64
}

type UserData struct {
	ID       uuid.UUID
	Username string
	Name     string
}
