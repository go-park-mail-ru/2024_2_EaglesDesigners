package model

// @Schema
type Chat struct {
	ChatId   int    `json:"chatId" example:"1"`
	ChatName string `json:"chatName" example:"Чат с пользователем 2"`
	// @Enum [personalMessages, group, channel]
	ChatType    string `json:"chatType" example:"personalMessages"`
	UsersId     []int  `json:"usersId" example:"1,2"`
	LastMessage string `json:"lastMessage" example:"Когда за кофе?"`
	AvatarURL   string `json:"avatarURL" example:"https://yandex-images.clstorage.net/bVLC53139/"`
}

type ChatsDTO struct {
	Chats []Chat `json:"chats"`
}
