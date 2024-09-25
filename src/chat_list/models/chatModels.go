package model

type Chat struct {
	ChatId   int    `json:"chatId"`
	ChatName string `json:"chatName"`
}


type ChatsDTO struct {
	Chats []Chat `json:"chats"`
}