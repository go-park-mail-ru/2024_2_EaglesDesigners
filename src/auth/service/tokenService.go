package service

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
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

func (s *TokenService) IsAuthorized(cookies []*http.Cookie) bool {
	token, err := parserCookies(cookies)
	if err != nil {
		return false
	}

	result, err := checkJWT(token)
	if err != nil {
		return false
	}

	return result
}

func (s *TokenService) CreateJWT(username string) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	id := s.userRepo.GetUserByUsername(username).ID

	payload := Payload{
		Sub:  id,
		Name: username,
		Iat:  time.Now().Unix(),
		Exp:  time.Now().Add(time.Hour * 72).Unix(),
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

func (s *TokenService) GetUserByJWT(cookies []*http.Cookie) model.User {
	token, err := parserCookies(cookies)
	if err != nil {
		return model.User{}
	}

	jwt := strings.Split(token, ".")
	if len(jwt) != 3 {
		return model.User{}
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(jwt[2])
	if err != nil {
		return model.User{}
	}

	var payload Payload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return model.User{}
	}

	return s.userRepo.GetUserByUsername(payload.Name)
}
