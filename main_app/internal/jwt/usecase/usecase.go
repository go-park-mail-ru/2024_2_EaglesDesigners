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

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"
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

func (u *Usecase) IsAuthorized(ctx context.Context, cookies []*http.Cookie) (user User, err error) {
	token, err := parseCookies(cookies)
	if err != nil {
		return user, err
	}

	result, err := checkJWT(token)
	if err != nil {
		return user, err
	}

	if !result {
		return user, errors.New("токен невалиден")
	}

	payload, err := getPayloadOfJWT(token)
	if err != nil {
		return user, err
	}

	user, err = u.GetUserByJWT(ctx, cookies)
	if err != nil {
		return user, err
	}

	if payload.Version != user.Version {
		return user, errors.New("токен устарел")
	}

	if payload.Exp < time.Now().Unix() {
		return user, errors.New("токен истек")
	}

	return user, nil
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

func (u *Usecase) GetUserByJWT(ctx context.Context, cookies []*http.Cookie) (User, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("запрошен поиск пользователь по jwt")

	token, err := parseCookies(cookies)
	if err != nil {
		return User{}, err
	}

	payload, err := getPayloadOfJWT(token)
	if err != nil {
		return User{}, err
	}

	log.Println("пользователь аутентификацирован")

	repoUser, err := u.repository.GetUserByUsername(ctx, payload.Sub)
	if err != nil {
		log.Errorf("пользователь не найден: %v", err)
		return User{}, err
	}

	user := convertToUser(repoUser)

	return user, nil
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

func checkJWT(token string) (bool, error) {
	jwt := strings.Split(token, ".")
	if len(jwt) != 3 {
		return false, errors.New("invalid token")
	}
	header := jwt[0]
	payload := jwt[1]
	signature := jwt[2]

	newToken, err := GeneratorJWT(header, payload, jwtSecret)
	if err != nil {
		return false, err
	}

	newSignature := strings.Split(newToken, ".")[2]

	return signature == newSignature, nil
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

func convertToUser(u auth.UserDAO) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Version:  u.Version,
	}
}