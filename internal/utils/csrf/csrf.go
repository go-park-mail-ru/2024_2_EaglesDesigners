package csrf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/google/uuid"
)

var jwtSecret = jwt.GenerateJWTSecret()

var errInvalidToken = errors.New("невалидный csrf токен")

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub string    `json:"sub"`
	ID  uuid.UUID `json:"id"`
	Exp int64     `json:"exp"`
}

func CreateCSRF(accessToken string) (string, error) {
	accessPayloadBase64 := strings.Split(accessToken, ".")[1]

	var accessPayload jwt.Payload

	accessPayloadBytes, err := base64.RawURLEncoding.DecodeString(accessPayloadBase64)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(accessPayloadBytes, &accessPayload)
	if err != nil {
		return "", err
	}

	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := Payload{
		Sub: accessPayload.Sub,
		ID:  accessPayload.ID,
		Exp: time.Now().Add(time.Hour * 24).Unix(),
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

	jwt, err := jwt.GeneratorJWT(headerEncoded, payloadEncoded, jwtSecret)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func CheckCSRF(token string, userID uuid.UUID, username string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errInvalidToken
	}

	headerBase64 := parts[0]
	payloadBase64 := parts[1]
	signatureBase64 := parts[2]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return err
	}

	var payload Payload

	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return err
	}

	if payload.Sub != username || payload.ID != userID {
		return errors.New("данные в токенах не совпадают")
	}

	if payload.Exp < time.Now().Unix() {
		return errors.New("токен истек")
	}

	newToken, err := jwt.GeneratorJWT(headerBase64, payloadBase64, jwtSecret)
	if err != nil {
		return err
	}

	newSignature := strings.Split(newToken, ".")[2]

	if signatureBase64 != newSignature {
		return errors.New("подпись не прошла проверку")
	}

	return nil
}
