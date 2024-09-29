package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/service"
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
	Message string           `json:"message" example:"Registration successful"`
	User    service.UserData `json:"user"`
}

// @Schema
type AuthResponse struct {
	User service.UserData `json:"user"`
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

func sendOKResponse(w http.ResponseWriter, message string) {
	response := SuccessResponse{Message: message}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
}
