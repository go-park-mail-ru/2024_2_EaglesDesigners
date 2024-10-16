package responser

import (
	"encoding/json"
	"log"
	"net/http"
)

// @Schema
type SuccessResponse struct {
	Message string `json:"message"`
}

// @Schema
type ErrorResponse struct {
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
