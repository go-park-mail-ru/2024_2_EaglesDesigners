package models

// @Schema
type UserDTO struct {
	Username     *string `json:"username" example:"user11"`
	Name         *string `json:"name" example:"Витек"`
	AvatarBase64 *string `json:"avatarBase64" example:"this is Base64 photo"`
}

// @Schema
type GetContactsResponseDTO struct {
	Contacts []User `json:"contacts"`
}

type User struct {
	Username     *string
	Name         *string
	AvatarBase64 *string
}

type UserDAO struct {
	Username  *string
	Name      *string
	AvatarURL *string
}
