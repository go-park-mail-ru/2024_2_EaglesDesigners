package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var jwtSecret = []byte("КТо пРочитАл тОт сдОхНет :)")

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	ID      int64  `json:"id"`
	Version int64  `json:"vrs"`
	Exp     int64  `json:"exp"`
}

type UserData struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func generatorJWT(header string, payload string) (string, error) {
	hmac := hmac.New(sha256.New, jwtSecret)
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

	newToken, err := generatorJWT(header, payload)
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
