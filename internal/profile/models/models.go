package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// @Schema
type UpdateProfileRequestDTO struct {
	ID        uuid.UUID       `json:"-"`
	Name      *string         `json:"name" example:"Vincent Vega"`
	Bio       *string         `json:"bio" example:"Не люблю сети"`
	Birthdate *time.Time      `json:"birthdate" example:"2024-04-13T08:30:00Z"`
	Avatar    *multipart.File `json:"-"`
}

// @Schema
type GetProfileResponseDTO struct {
	Name      *string    `json:"name" example:"Vincent Vega"`
	Bio       *string    `json:"bio" example:"Не люблю сети"`
	AvatarURL *string    `json:"avatarURL" example:"/2024_2_eaglesDesigners/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png"`
	Birthdate *time.Time `json:"birthdate" example:"2024-04-13T08:30:00Z"`
}

type Profile struct {
	ID        uuid.UUID
	Name      *string
	Bio       *string
	Avatar    *multipart.File
	Birthdate *time.Time
}

type ProfileData struct {
	Name       *string
	Bio        *string
	AvatarPath *string
	Birthdate  *time.Time
}

type ProfileDAO struct {
	ID         uuid.UUID
	Name       *string
	Bio        *string
	AvatarPath *string
	Birthdate  *time.Time
}

type ProfileDataDAO struct {
	Name       *string
	Bio        *string
	AvatarPath *string
	Birthdate  *time.Time
}
