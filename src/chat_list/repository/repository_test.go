package repository

import (
	"testing"

	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
	"github.com/stretchr/testify/require"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
)

func TestGetUser1Chats(t *testing.T) {
	chats := NewChatRepository().GetUserChats(&model.User{ID: 1})
	require.Equal(t, chats, []chatModel.Chat{
		chatModel.Chat{
			ChatId:      1,
			ChatName:    "Чат с пользователем 2",
			ChatType:    "personalMessages",
			UsersId:     []int{1, 2},
			LastMessage: "Когда за кофе?",
			AvatarURL:   "https://yandex-images.clstorage.net/bVLC53139/667e899dbzgI/Lec3og97oM2J8jgAbwmbs1UEQ_j2WQe6H7Tz0tGHlNUDiLp06xNO9LooehtZCLyucrVfOV3bNS1vNvr_fMoMLbniE8frC6CczUKcwc_ImueU0HKs18lHz490gERWwAOWtD4IttmRuiGPuG9PrfwYeJTUCT5PeyM6mMdYuXvvucreJwTwaprjvy1RSHHf2XlUxVagbjT_Z3s54KP1tiFyt1ZNSQbbE3rzTsqefIsOGIsUYXo-bNgSKZq1WSlJDWhoz9XEo5uL4K6ts3gATet5lVTXMkxWjOY7KIFDVJZywRQU_5pGKzNLwd44T06oTSlJF8LpfEw5cni4VlkKqL3pur7HtGJr6fJM_7N44i2JKwV39mKvRF_0WggAkvaH0qM3dF3rRmhTySM-KU0umyyYWiVw2A7POeGIuvaYeurMKuh91cfSm1uT3gyRmwOOe-pFVEUzzFSupwioEVOHBbBzZ2Tv-6XKcnrQn8kvjEot-5skUlksz1lzKIsGKytJHes6zpeFgbqbE_-fI5hA_TpqVeemkPzm7qcL-UCy5NVyUTYUrfi2K2JIwT85fHzLzVhIFiBIH-9b0Xt45-m5aDwpWu_VtENpCtPcbhE7gL9aisTlNiP8hE2kW1qwEiTkwILXxe9adbgRuvDt6U59ag4qWbVyOm7eG4Nq2QW5mTveSXpcdzTz-TmS_F8iiCEc6Or0JvZy3ldd9-q7MMEWl-FS5YT_GUTZQQqiztuNzust-xr0cJjMzfuSmqqlyvuIPKlJbadlI-l7IS6vUWixLdt6p1dG0dy3bMdYunPBRcWioYYX3BsV2fP7oq26jA6YvgmLxQKqnr_IU0tJJ4greQw5Sp0HhdFYu4DeXGFbEF3JOWb3ZIJvR851WQnhM2UX8VKF983KBAoDOOK9qp1_aC9YK0WCa_8P-UGZ2Eaqy5s9C4qsNofD6Rkg3qyRGCPveNplxMSyzjbctEk5w6C2k",
		},
		chatModel.Chat{
			ChatId:      2,
			ChatName:    "МГТУ",
			ChatType:    "group",
			UsersId:     []int{1, 2, 3},
			LastMessage: "У нас еще вечером одна пара",
			AvatarURL:   "https://polymerbranch.com/wp-content/uploads/2023/11/mgtu.webp",
		},
		chatModel.Chat{
			ChatId:      3,
			ChatName:    "Смешные картинки",
			ChatType:    "channel",
			UsersId:     []int{1, 2, 3},
			LastMessage: "Это была не смешная картинка",
			AvatarURL:   "https://i.pinimg.com/736x/dd/39/3f/dd393f8f2f8293c670f1fb4a74abef94.jpg",
		},
		chatModel.Chat{
			ChatId:      4,
			ChatName:    "Технопарк",
			ChatType:    "group",
			UsersId:     []int{1, 2, 3, 4},
			LastMessage: "Надо фронт с беком соединить",
			AvatarURL:   "https://i.pinimg.com/originals/0e/7e/c4/0e7ec411745a33daafa689df1d207667.jpg",
		},
	})
}

func TestGetNothing(t *testing.T) {
	chats := NewChatRepository().GetUserChats(&model.User{ID: 5})
	require.Equal(t, chats, []chatModel.Chat{})
}