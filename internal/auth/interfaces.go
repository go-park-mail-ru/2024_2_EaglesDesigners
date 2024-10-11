package auth

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/utils"
)

type AuthService interface {
	Authenticate(username, password string) bool
	Registration(username, name, password string) error
	GetUserDataByUsername(username string) (utils.UserData, error)
}

type TokenService interface {
	CreateJWT(username string) (string, error)
	GetUserDataByJWT(cookies []*http.Cookie) (utils.UserData, error)
	GetUserByJWT(cookies []*http.Cookie) (model.User, error)
	IsAuthorized(cookies []*http.Cookie) error
}

type UserRepository interface {
	GetUserByUsername(username string) (model.User, error)
	CreateUser(username, name, password string) error
}
