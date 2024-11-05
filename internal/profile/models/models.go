package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// @Schema
type UpdateProfileRequestDTO struct {
	ID        uuid.UUID       `json:"-" valid:"-"`
	Name      *string         `json:"name" example:"Vincent Vega" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	Bio       *string         `json:"bio" example:"Не люблю сети" valid:"optional"`
	Birthdate *time.Time      `json:"birthdate" example:"2024-04-13T08:30:00Z" valid:"optional"`
	Avatar    *multipart.File `json:"-" valid:"-"`
}

// @Schema
type GetProfileResponseDTO struct {
	Name      *string    `json:"name" example:"Vincent Vega" valid:"matches(^[а-яА-Яa-zA-Z0-9_ ]+$)"`
	Bio       *string    `json:"bio" example:"Не люблю сети" valid:"optional"`
	AvatarURL *string    `json:"avatarURL" example:"/uploads/avatar/f0364477-bfd4-496d-b639-d825b009d509.png" valid:"matches(^/uploads/avatar/[a-zA-Z0-9\\-]+\\.png$),optional"`
	Birthdate *time.Time `json:"birthdate" example:"2024-04-13T08:30:00Z" valid:"optional"`
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
