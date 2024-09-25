package repository

import (
	"log"

	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
)


var keys = map[int64][]chatModel.Chat {
	0: []chatModel.Chat{
		chatModel.Chat{
			ChatId: 1,
			ChatName: "1",
			ChatType: "personalMessages",
			UsersId: []int{1, 2},
		},
		chatModel.Chat{
			ChatId: 2,
			ChatName: "2",
			ChatType: "group",
			UsersId: []int{1, 2, 3},
		},
		chatModel.Chat{
			ChatId: 3,
			ChatName: "3",
			ChatType: "channel",
			UsersId: []int{1, 2, 3},
		},
	},
	2:  []chatModel.Chat{
		chatModel.Chat{
			ChatId: 1,
			ChatName: "1",
			ChatType: "personalMessages",
			UsersId: []int{1, 2},
		},
	},
	3:  []chatModel.Chat{
		chatModel.Chat{
			ChatId: 2,
			ChatName: "2",
			ChatType: "group",
			UsersId: []int{1, 2, 3},
		},
		chatModel.Chat{
			ChatId: 4,
			ChatName: "4",
			ChatType: "group",
			UsersId: []int{1, 2, 3, 5},
		},
	},
	4:  []chatModel.Chat{
		chatModel.Chat{
			ChatId: 5,
			ChatName: "5",
			ChatType: "channel",
			UsersId: []int{4},
		},
		chatModel.Chat{
			ChatId: 4,
			ChatName: "4",
			ChatType: "group",
			UsersId: []int{1, 2, 3, 5},
		},
	},
}

func GetUserChats(user *userModel.User) []chatModel.Chat {
	chats, ok := keys[user.ID]
	log.Println(chats)
	
	if !ok {
		return []chatModel.Chat{}
	}
	
	return chats
}

