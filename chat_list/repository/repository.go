package repository

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Chat struct {
	ChatId   int    `json:"chatId"`
	ChatName string `json:"chatName"`
}

func GetUserChats(user *User) []Chat {

	return nil
}
