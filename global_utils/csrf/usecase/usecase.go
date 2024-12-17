package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/csrf/models"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mocks.go

var jwtSecret = GenerateJWTSecret()

type repository interface {
	GetUserByUsername(ctx context.Context, username string) (auth.UserDAO, error)
	CreateUser(ctx context.Context, username, name, password string) error
}

type Usecase struct {
	repository repository
}

func NewUsecase(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (u *Usecase) CreateJWT(ctx context.Context, username string) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	user, err := u.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	payload := Payload{
		Sub:     user.Username,
		Name:    user.Name,
		ID:      user.ID,
		Version: user.Version,
		Exp:     time.Now().Add(time.Hour * 24).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	jwt, err := GeneratorJWT(headerEncoded, payloadEncoded, jwtSecret)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (u *Usecase) GetUserDataByJWT(cookies []*http.Cookie) (UserData, error) {
	token, err := parseCookies(cookies)
	if err != nil {
		return UserData{}, err
	}
	payload, err := getPayloadOfJWT(token)
	if err != nil {
		return UserData{}, err
	}

	data := UserData{
		ID:       payload.ID,
		Username: payload.Sub,
		Name:     payload.Name,
	}

	return data, nil
}

func GenerateJWTSecret() []byte {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		log.Fatalf("Ошибка при генерации jwtSecret: %v", err)
	}
	return secret
}

func GeneratorJWT(header string, payload string, secret []byte) (string, error) {
	hmac := hmac.New(sha256.New, secret)
	hmac.Write([]byte(header + "." + payload))
	signature := hmac.Sum(nil)

	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	jwt := header + "." + payload + "." + signatureEncoded

	return jwt, nil
}

func parseCookies(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie does not exist")
}

func getPayloadOfJWT(token string) (payload Payload, err error) {
	jwt := strings.Split(token, ".")

	if len(jwt) != 3 {
		return payload, errors.New("невалидный jwt token")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(jwt[1])
	if err != nil {
		return payload, errors.New("невалидный jwt token")
	}

	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return payload, errors.New("невалидный jwt token")
	}

	return payload, nil
}
