package responser

import (
	"encoding/json"
	"log"
	"net/http"
)

// @Schema
type SuccessResponse struct {
	Message string `json:"message" example:"success message"`
}

// @Schema
type ErrorResponse struct {
	Error  string `json:"error" example:"error message"`
	Status string `json:"status" example:"error"`
}

func SendOK(w http.ResponseWriter, message string, statusCode int) {
	response := SuccessResponse{Message: message}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, errorMessage string, statusCode int) {
	log.Printf("Отправлен код %d. ОШИБКА: %s \n", statusCode, errorMessage)

	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	SendError(w, "Method not allowed", http.StatusUnauthorized)
}

// SendStruct отправляет полученный экземпляр структуры в формате json с статусом кода statusCode.
func 	SendStruct(w http.ResponseWriter, response any, statusCode int) {
	jsonResp, err := json.Marshal(response)
	if err != nil {
		SendError(w, "Failed to create response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}
