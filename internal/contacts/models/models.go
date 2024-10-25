package models

import "github.com/google/uuid"

// @Schema
type ContactDTO struct {
	ID       string `json:"id" example:"08a0f350-e122-467b-8ba8-524d2478b56e"`
	Username string `json:"username" example:"user11"`

	// can be nil
	Name         *string `json:"name" example:"Витек"`
	AvatarBase64 *string `json:"avatarBase64" example:"this is Base64 photo"`
}

// @Schema
type AddContactReqDTO struct {
	Username string `json:"contactUsername" example:"user11"`
}

// @Schema
type GetContactsRespDTO struct {
	Contacts []ContactDTO `json:"contacts"`
}

type Contact struct {
	ID           string
	Username     string
	Name         *string
	AvatarBase64 *string
}

type ContactData struct {
	UserID          string
	ContactUsername string
}

type ContactDAO struct {
	ID        uuid.UUID
	Username  string
	Name      *string
	AvatarURL *string
}

type ContactDataDAO struct {
	UserID          string
	ContactUsername string
}
