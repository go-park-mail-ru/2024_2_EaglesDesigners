package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.octolab.org/pointer"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/models"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mocks.go

func init() {
	prometheus.MustRegister(newUserMetric)
}

type repository interface {
	GetUserByUsername(ctx context.Context, username string) (models.UserDAO, error)
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

func (u *Usecase) Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	user, err := u.repository.GetUserByUsername(ctx, in.GetUsername())
	if err != nil {
		return nil, err
	}

	return &authv1.AuthResponse{
		IsAuthenticated: DoPasswordsMatch(user.Password, in.GetPassword()),
	}, nil
}

var newUserMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "count_of_registered_users",
		Help: "countOfHits",
	},
	nil, // no labels for this metric
)

func (u *Usecase) Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	if len(in.GetUsername()) < 6 || len(in.GetPassword()) < 8 || len(in.GetName()) < 1 {
		log.Println("не удалось создать юзера: данные не прошли валидацию")
		return nil, errors.New("bad data")
	}

	hashed := HashPassword(in.GetPassword())
	err := u.repository.CreateUser(ctx, in.GetUsername(), in.GetName(), hashed)
	if err != nil {
		log.Println("не удалось создать юзера: ", err)
		return nil, err
	}

	log.Println("пользователь создан")
	metric.IncMetric(*newUserMetric)
	return &authv1.Nothing{Dummy: true}, nil
}

func (u *Usecase) GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	user, err := u.repository.GetUserByUsername(ctx, in.GetUsername())
	if err != nil {
		return nil, err
	}

	log.Println("пользователь получен")

	return &authv1.GetUserDataByUsernameResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Name:      user.Name,
		AvatarURL: pointer.ValueOfString(user.AvatarURL),
	}, nil
}

func getSalt() []byte {
	return []byte{93, 108, 25, 43, 92, 102, 255, 179, 11, 87, 186, 198, 254, 160, 164, 56}
}

func HashPassword(password string) string {
	passwordBytes := []byte(password)
	sha512Hasher := sha512.New()
	passwordBytes = append(passwordBytes, getSalt()...)
	sha512Hasher.Write(passwordBytes)
	hashedPasswordBytes := sha512Hasher.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	currPasswordHash := HashPassword(currPassword)
	return hashedPassword == currPasswordHash
}

// JWT

var jwtSecret = GenerateJWTSecret()

func (u *Usecase) IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	result, err := checkJWT(in.GetToken())
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, errors.New("токен невалиден")
	}

	payload, err := getPayloadOfJWT(in.GetToken())
	if err != nil {
		return nil, err
	}

	user, err := u.GetUserByJWT(ctx, in)
	if err != nil {
		return nil, err
	}

	if payload.Version != int64(user.Version) {
		return user, errors.New("токен устарел")
	}

	if payload.Exp < time.Now().Unix() {
		return user, errors.New("токен истек")
	}

	return user, nil
}

func (u *Usecase) CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error) {
	header := models.Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	user, err := u.repository.GetUserByUsername(ctx, in.GetUsername())
	if err != nil {
		return nil, err
	}

	payload := models.Payload{
		Sub:     user.Username,
		Name:    user.Name,
		ID:      user.ID,
		Version: user.Version,
		Exp:     time.Now().Add(time.Hour * 24).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	jwt, err := GeneratorJWT(headerEncoded, payloadEncoded, jwtSecret)
	if err != nil {
		return nil, err
	}

	return &authv1.Token{Token: jwt}, nil
}

func (u *Usecase) GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Debugln("запрошен поиск пользователь по jwt")

	payload, err := getPayloadOfJWT(in.GetToken())
	if err != nil {
		return nil, err
	}

	log.Debugln("пользователь аутентификацирован")

	repoUser, err := u.repository.GetUserByUsername(ctx, payload.Sub)
	if err != nil {
		log.Errorf("пользователь не найден: %v", err)
		return nil, err
	}

	user := convertToUser(repoUser)

	return &authv1.UserJWT{
		ID:       user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Password: user.Password,
		Version:  user.Version,
	}, nil
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

func getPayloadOfJWT(token string) (payload models.Payload, err error) {
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

func convertToUser(u models.UserDAO) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
		Version:  u.Version,
	}
}
