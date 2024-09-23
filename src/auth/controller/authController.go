package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
)

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

func (c *AuthController) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
		return
	}

	var creds AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "invalid format JSON", http.StatusBadRequest)
		return
	}

	if c.authService.Authenticate(creds.Username, creds.Password) {
		token, err := c.tokenService.CreateJWT(creds.Username)
		if err != nil {
			sendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
			return
		}

		cookie := &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(48 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, cookie)

		sendOKResponse(w, "Authentication successful")
	} else {
		sendErrorResponse(w, "Incorrect login or password", http.StatusUnauthorized)
	}
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
		return
	}

	var creds RegisterCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		sendErrorResponse(w, "Invalid format JSON", http.StatusBadRequest)
		return
	}

	if err := c.authService.Registation(creds.Username, creds.Password); err != nil {
		sendErrorResponse(w, "A user with that username already exists", http.StatusConflict)
	} else {
		sendOKResponse(w, "Registration successful")
	}
}

func sendOKResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{
		"error":  message,
		"status": "error",
	}

	json.NewEncoder(w).Encode(response)
}
