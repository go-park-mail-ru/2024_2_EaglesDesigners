package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"net/http"

	repo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type repository interface {
	GetUserByUsername(ctx context.Context, username string) (repo.User, error)
	CreateUser(ctx context.Context, username, name, password string) error
}

type token interface {
	CreateJWT(ctx context.Context, username string) (string, error)
	GetUserDataByJWT(cookies []*http.Cookie) (jwt.UserData, error)
	GetUserByJWT(ctx context.Context, cookies []*http.Cookie) (jwt.User, error)
}

type Usecase struct {
	repository repository
	token      token
}

func NewUsecase(repository repository, token token) *Usecase {
	return &Usecase{
		repository: repository,
		token:      token,
	}
}

func (u *Usecase) Authenticate(ctx context.Context, username, password string) bool {
	user, err := u.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return false
	}
	return DoPasswordsMatch(user.Password, password)
}

func (u *Usecase) Registration(ctx context.Context, username, name, password string) error {
	if len(username) < 6 || len(password) < 8 || len(name) < 1 {
		return errors.New("bad data")
	}

	hashed := HashPassword(password)
	err := u.repository.CreateUser(ctx, username, name, hashed)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) GetUserDataByUsername(ctx context.Context, username string) (UserData, error) {
	user, err := u.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return UserData{}, err
	}

	userData := UserData{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	return userData, nil
}

func getSalt() []byte {
	return []byte{93, 108, 25, 43, 92, 102, 255, 179, 11, 87, 186, 198, 254, 160, 164, 56}
}

func HashPassword(password string) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, getSalt()...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	var currPasswordHash = HashPassword(currPassword)
	return hashedPassword == currPasswordHash
}
