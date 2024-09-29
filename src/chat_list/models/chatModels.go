package model

type Chat struct {
	ChatId   int    `json:"chatId"`
	ChatName string `json:"chatName"`
	ChatType string `json:"chatType"`
	UsersId []int	`json:"usersId"`
	LastMessage string `json:"lastMessage"`
	AvatarURL string `json:"avatarURL"`
}

type ChatsDTO struct {
	Chats []Chat `json:"chats"`
}