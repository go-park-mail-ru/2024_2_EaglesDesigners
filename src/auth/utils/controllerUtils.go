package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// AuthCredentials represents the credentials required for authentication
// @Schema
type AuthCredentials struct {
	Username string `json:"username" example:"user1"`
	Password string `json:"password"  example:"pass1"`
}

// RegisterCredentials represents the credentials required for user registration
// @Schema
type RegisterCredentials struct {
	Username string `json:"username" example:"killer1994"`
	Name     string `json:"name" example:"Vincent Vega"`
	Password string `json:"password" example:"go_do_a_crime"`
}

// @Schema
type RegisterResponse struct {
	Message string   `json:"message" example:"Registration successful"`
	User    UserData `json:"user"`
}

// @Schema
type AuthResponse struct {
	User UserData `json:"user"`
}

// @Schema
type SuccessResponse struct {
	Message string `json:"message"`
}

// @Schema
type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}

type SignupResponse struct {
	Error  string `json:"error"`
	Status string `json:"status" example:"error"`
}

func SendOKResponse(w http.ResponseWriter, message string, statusCode int) {
	response := SuccessResponse{Message: message}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	log.Printf("Отправлен код %d. ОШИБКА: %s \n", statusCode, errorMessage)

	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	SendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
}
