package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	model "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/models"
	"github.com/stretchr/testify/require"
)

type tokenService struct {
	User model.User
	err  error
}

func (m *tokenService) GetUserByJWT(cookies []*http.Cookie) (repository.User, error) {
	return m.User, m.err
}

func (m *tokenService) CreateJWT(username string) (string, error) {
	return "", nil
}

func (m *tokenService) GetUserDataByJWT(cookies []*http.Cookie) (model.UserData, error) {
	return utils.UserData{}, nil
}

func (m *tokenService) IsAuthorized(cookies []*http.Cookie) error {
	return nil
}

type chatRepositoryMoc struct {
	chats []chatModel.Chat
}

func (r *chatRepositoryMoc) GetUserChats(user *model.User) []chatModel.Chat {
	return r.chats
}

func TestSuccess(t *testing.T) {
	ts := &tokenService{
		User: model.User{
			ID: 1,
		},
		err: nil,
	}
	chats := []chatModel.Chat{
		chatModel.Chat{
			ChatId:      1,
			ChatName:    "Чат с пользователем 1",
			ChatType:    "personalMessages",
			UsersId:     []int{1, 2},
			LastMessage: "Когда за кофе?",
			AvatarURL:   "https://get.wallhere.com/photo/2048x1280-px-animals-baby-cat-cats-cute-fat-fluffy-grass-grey-kitten-kittens-1913313.jpg",
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
	}
	cr := &chatRepositoryMoc{
		chats: chats,
	}

	serv := NewChatService(ts, cr)
	chats, err := serv.GetChats([]*http.Cookie{})

	if err != nil {
		t.Errorf("Expected no errors, got %s", err)
	}
	require.Equal(t, chats, chats, "The two words should be the same.")

}

func TestUnauthorized(t *testing.T) {
	ts := &tokenService{
		User: model.User{},
		err:  errors.New("Unauthorized"),
	}
	cr := &chatRepositoryMoc{}

	serv := NewChatService(ts, cr)

	_, err := serv.GetChats([]*http.Cookie{})

	if err == nil {
		t.Fail()
	}
}
