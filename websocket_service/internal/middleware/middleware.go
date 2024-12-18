package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/csrf"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/authv1"
)

type Delivery struct {
	client authv1.AuthClient
}

func New(client authv1.AuthClient) *Delivery {
	return &Delivery{
		client: client,
	}
}

type contextKey string

const (
	UserIDKey    contextKey = "userId"
	UserKey      contextKey = "user"
	MuxParamsKey contextKey = "muxParams"
)

var errTokenExpired = errors.New("токен истек")

func (d *Delivery) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := d.parseCookies(r.Cookies())
		if err != nil {
			log.Println("не получилось получить токен")
			responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		user, err := d.client.IsAuthorized(ctx, &authv1.Token{Token: token})
		if err == errTokenExpired {
			log.Println("токен истек, создается новый токен")
			d.setTokens(w, r, user.Username)
		}

		if err != nil && err != errTokenExpired {
			log.Println(err)
			log.Println("токен невалиден")
			responser.SendError(ctx, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, UserKey, convertFromGRPCUser(user))
		ctx = context.WithValue(ctx, MuxParamsKey, mux.Vars(r))

		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (d *Delivery) setTokens(w http.ResponseWriter, r *http.Request, username string) error {
	ctx := r.Context()
	grcpResp, err := d.client.CreateJWT(ctx, &authv1.CreateJWTRequest{Username: username})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    grcpResp.GetToken(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	csrf, err := csrf.CreateCSRF(grcpResp.GetToken())
	if err != nil {
		log.Printf("Auth setTokens: не удалось создать csrf токен: %v", err)
		return err
	}

	w.Header().Set("X-CSRF-Token", csrf)
	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

	return nil
}

func (d *Delivery) parseCookies(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie does not exist")
}

type User struct {
	ID       uuid.UUID
	Username string
	Name     string
	Password string
	Version  int64
}

func convertFromGRPCUser(user *authv1.UserJWT) User {
	return User{
		ID:       uuid.MustParse(user.GetID()),
		Username: user.GetUsername(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
		Version:  user.GetVersion(),
	}
}
