package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/model"
	"github.com/stretchr/testify/require"
)

var usersT = map[string]model.User{
	"user11": {
		ID:       1,
		Username: "user11",
		Name:     "Бал Матье",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user22": {
		ID:       2,
		Username: "user22",
		Name:     "Жабка Пепе",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user33": {
		ID:       3,
		Username: "user33",
		Name:     "Dr Peper",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user44": {
		ID:       4,
		Username: "user44",
		Name:     "Vincent Vega",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
}

func TestGetAllStoredUsers(t *testing.T) {
	rep := NewUserRepository()
	for _, usr := range usersT {
		userFromStorage, err := rep.GetUserByUsername(usr.Username)

		if err != nil {
			t.Fail()
		}

		require.Equal(t, userFromStorage, usr)
	}
}

func TestGetUnstoredUser(t *testing.T) {
	rep := NewUserRepository()
	username := "aboba"

	_, err := rep.GetUserByUsername(username)

	if err == nil {
		t.Fail()
	}

}

func TestAddNewUser(t *testing.T) {
	username := "aboba"
	name := "Oleg Kizaru"
	password := "12345678"

	rep := NewUserRepository()
	err := rep.CreateUser(username, name, password)

	if err != nil {
		t.Fail()
	}
	newUsr, err := rep.GetUserByUsername("aboba")

	if err != nil {
		t.Fail()
	}
	require.Equal(t, username, newUsr.Username)
	require.Equal(t, password, newUsr.Password)
	require.Equal(t, name, newUsr.Name)
}

func TestAddDuplicateUserFail(t *testing.T) {
	username := "abobus"
	name := "Oleg Kizaru"
	password := "12345678"

	rep := NewUserRepository()
	err := rep.CreateUser(username, name, password)
	if err != nil {
		t.Fail()
	}

	err = rep.CreateUser(username, name, password)
	if err == nil {
		t.Fail()
	}
}
