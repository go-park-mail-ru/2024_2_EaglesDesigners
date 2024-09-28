package controller

import (
	"encoding/json"
	"net/http"
)

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterCredentials struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
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

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Method not allowed", http.StatusUnauthorized)
}
