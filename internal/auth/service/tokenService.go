package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/utils"
)

type TokenService struct {
	userRepo auth.UserRepository
}

func NewTokenService(userRepo auth.UserRepository) *TokenService {
	return &TokenService{
		userRepo: userRepo,
	}
}

func (s *TokenService) IsAuthorized(cookies []*http.Cookie) error {
	token, err := utils.ParseCookies(cookies)
	if err != nil {
		return err
	}

	result, err := utils.CheckJWT(token)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("invalid token")
	}

	payload, err := utils.GetPayloadOfJWT(token)
	if err != nil {
		return err
	}

	user, err := s.GetUserByJWT(cookies)
	if err != nil {
		return err
	}

	if payload.Version != user.Version {
		return errors.New("token outdated")
	}

	if payload.Exp < time.Now().Unix() {
		return errors.New("token expired")
	}

	return nil
}

func (s *TokenService) CreateJWT(username string) (string, error) {
	header := utils.Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	payload := utils.Payload{
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

	jwt, err := utils.GeneratorJWT(headerEncoded, payloadEncoded)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *TokenService) GetUserByJWT(cookies []*http.Cookie) (model.User, error) {
	log.Println("Запрошен поиск пользователь по jwt")

	token, err := utils.ParseCookies(cookies)
	if err != nil {
		return model.User{}, err
	}

	payload, err := utils.GetPayloadOfJWT(token)
	if err != nil {
		return model.User{}, err
	}

	log.Println("Пользователь аутентификацирован")

	return s.userRepo.GetUserByUsername(payload.Sub)
}

func (s *TokenService) GetUserDataByJWT(cookies []*http.Cookie) (utils.UserData, error) {
	token, err := utils.ParseCookies(cookies)
	if err != nil {
		return utils.UserData{}, err
	}
	payload, err := utils.GetPayloadOfJWT(token)
	if err != nil {
		return utils.UserData{}, err
	}

	data := utils.UserData{
		ID:       payload.ID,
		Username: payload.Sub,
		Name:     payload.Name,
	}

	return data, nil
}
