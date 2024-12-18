package responser

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
)

const (
	InvalidJSONError  = "invalid format JSON"
	ResponseError     = "failed to create response"
	UserNotFoundError = "user not found"
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

func SendError(ctx context.Context, w http.ResponseWriter, errorMessage string, statusCode int) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Errorf("Отправлен код %d. ОШИБКА: %s", statusCode, errorMessage)

	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	pc, _, _, _ := runtime.Caller(1)
	funcPath := runtime.FuncForPC(pc).Name()

	metric.PushError(funcPath, statusCode)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusUnauthorized
	errorMessage := "Method not allowed"

	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Errorf("Отправлен код %d. ОШИБКА: %s", statusCode, errorMessage)
}

// SendStruct отправляет полученный экземпляр структуры в формате json с статусом кода statusCode.
func SendStruct(ctx context.Context, w http.ResponseWriter, response interface{}, statusCode int) {
	jsonResp, err := json.Marshal(response)
	if err != nil {
		SendError(ctx, w, "Failed to create response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}

// SendStruct отправляет полученный экземпляр структуры в формате json с статусом кода statusCode.
func SendJson(ctx context.Context, w http.ResponseWriter, jsonResp []byte, err error, statusCode int) {
	if err != nil {
		SendError(ctx, w, "Failed to create response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}
