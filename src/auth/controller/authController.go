package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
)

type AuthController struct {
	authService  service.AuthService
	tokenService service.TokenService
}

func NewAuthController(authService service.AuthService, tokenService service.TokenService) *AuthController {
	return &AuthController{
		authService:  authService,
		tokenService: tokenService,
	}
}

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if c.authService.Authenticate(creds.Username, creds.Password) {
		err := c.setToken(w, creds.Username)
		if err != nil {
			sendErrorResponse(w, "Invalid format JSON", http.StatusUnauthorized)
			return
		}

		sendOKResponse(w, "Authentication successful")
	} else {
		sendErrorResponse(w, "Incorrect login or password", http.StatusUnauthorized)
	}
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds RegisterCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" || creds.Name == "" {
		sendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if err := c.authService.Registation(creds.Username, creds.Name, creds.Password); err != nil {
		sendErrorResponse(w, "A user with that username already exists", http.StatusConflict)
	} else {
		sendOKResponse(w, "Registration successful")
	}
}

func (c *AuthController) AuthHandler(w http.ResponseWriter, r *http.Request) {
	data, err := c.tokenService.GetUserDataByJWT(r.Cookies())
	log.Println("/auth cookie: ", data)
	if err != nil {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := struct {
		User service.UserData `json:"user"`
	}{
		User: data,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (c *AuthController) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.isAuthorized(w, r)
		if err == errors.New("token expired") {
			log.Println("token expired, create new token")
			user, err := c.tokenService.GetUserByJWT(r.Cookies())
			if err != nil {
				sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			c.setToken(w, user.Username)
		}
		if err != nil {
			sendErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (c *AuthController) isAuthorized(w http.ResponseWriter, r *http.Request) error {
	err := c.tokenService.IsAuthorized(r.Cookies())
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthController) setToken(w http.ResponseWriter, username string) error {
	token, err := c.tokenService.CreateJWT(username)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, cookie)
	return nil
}
