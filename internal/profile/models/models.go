package models

import (
	"database/sql"
	"time"
)

// @Schema
type UpdateProfileRequestDTO struct {
	Username     string
	Name         *string    `json:"name" example:"Vincent Vega"`
	Bio          *string    `json:"bio" example:"Не люблю сети"`
	AvatarBase64 *string    `json:"avatarBase64" example:"this is Base64 photo"`
	Birthdate    *time.Time `json:"birthdate" example:"2024-04-13T08:30:00Z"`
}

// @Schema
type GetProfileRequestDTO struct {
	Username string `json:"username" example:"killer1994"`
}

// @Schema
type GetProfileResponseDTO struct {
	Bio          *string    `json:"bio" example:"Не люблю сети"`
	AvatarBase64 *string    `json:"avatarBase64" example:"this is Base64 photo"`
	Birthdate    *time.Time `json:"birthdate" example:"2024-04-13T08:30:00Z"`
}

type Profile struct {
	Username     string
	Name         *string
	Bio          *string
	AvatarBase64 *string
	Birthdate    *time.Time
}

type ProfileData struct {
	Bio       *string
	AvatarURL *string
	Birthdate *time.Time
}

type ProfileDAO struct {
	Username  string
	Name      sql.NullString
	Bio       sql.NullString
	AvatarURL sql.NullString
	Birthdate sql.NullTime
}

type ProfileDataDAO struct {
	Bio       sql.NullString
	AvatarURL sql.NullString
	Birthdate sql.NullTime
}
