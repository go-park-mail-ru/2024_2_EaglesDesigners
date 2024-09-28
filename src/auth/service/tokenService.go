package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/model"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/repository"
)

type TokenService struct {
	userRepo repository.UserRepository
}

func NewTokenService(userRepo repository.UserRepository) *TokenService {
	return &TokenService{
		userRepo: userRepo,
	}
}

func (s *TokenService) IsAuthorized(cookies []*http.Cookie) error {
	token, err := parseCookies(cookies)
	if err != nil {
		return err
	}

	result, err := checkJWT(token)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("invalid token")
	}

	payload, err := getPayloadOfJWT(token)
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
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	payload := Payload{
		Sub:     user.Username,
		Name:    user.Name,
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

	jwt, err := generatorJWT(headerEncoded, payloadEncoded)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *TokenService) GetUserByJWT(cookies []*http.Cookie) (model.User, error) {
	token, err := parseCookies(cookies)
	if err != nil {
		return model.User{}, err
	}

	payload, err := getPayloadOfJWT(token)
	if err != nil {
		return model.User{}, err
	}

	return s.userRepo.GetUserByUsername(payload.Sub)
}

func (s *TokenService) GetUserDataByJWT(cookies []*http.Cookie) (UserData, error) {
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
