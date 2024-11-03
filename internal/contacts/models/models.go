package models

import "github.com/google/uuid"

// @Schema
type ContactRespDTO struct {
	ID       string `json:"id" example:"08a0f350-e122-467b-8ba8-524d2478b56e"`
	Username string `json:"username" example:"user11"`

	// can be nil
	Name      *string `json:"name" example:"Витек"`
	AvatarURL *string `json:"avatarURL" example:"/uploads/avatar/642c5a57-ebc7-49d0-ac2d-f2f1f474bee7.png"`
}

// @Schema
type ContactReqDTO struct {
	Username string `json:"contactUsername" example:"user11"`
}

// @Schema
type GetContactsRespDTO struct {
	Contacts []ContactRespDTO `json:"contacts"`
}

type Contact struct {
	ID        string
	Username  string
	Name      *string
	AvatarURL *string
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
