package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	websocketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/usecase"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5024,
	WriteBufferSize: 5024,

	// CheckOrigin: func(r *http.Request) bool {
	// 	allowedOrigins := []string{
	// 		"http://127.0.0.1:8001",
	// 		"https://127.0.0.1:8001",
	// 		"http://localhost:8001",
	// 		"https://localhost:8001",
	// 		"http://213.87.152.18:8001",
	// 		"http://212.233.98.59:8001",
	// 		"https://213.87.152.18:8001",
	// 		"http://212.233.98.59:8080",
	// 		"https://212.233.98.59:8080",
	// 	}

	// 	for _, origin := range allowedOrigins {
	// 		if r.Header.Get("Origin") == origin {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// },
}

type Webcosket struct {
	usecase websocketUsecase.WebsocketUsecase
}

func NewWebsocket(usecase websocketUsecase.WebsocketUsecase) Webcosket {
	return Webcosket{
		usecase: usecase,
	}
}

func (h *Webcosket) HandleConnection(w http.ResponseWriter, r *http.Request) {
	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	// начало

	userId, _ := uuid.Parse("39a9aea0-d461-437d-b4eb-bf030a0efc80")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer log.Println("Message delivery: websocket is closing")
	defer conn.Close()

	eventChannel := make(chan websocketUsecase.AnyEvent, 10)

	err = h.usecase.InitBrokersForUser(userId, eventChannel)
	if err != nil {
		log.Errorf("Не удалось иницировать брокеры для пользователя")
		SendError(r.Context(), w, "Нет нужных параметров", http.StatusInternalServerError)

		return
	}

	// пока соеденено
	duration := 500 * time.Millisecond

	for {
		select {
		case message := <-eventChannel:
			// запись новых сообщений
			log.Println("Message delivery websocket: получены новые сообщения")

			conn.WriteJSON(message.Event)

		default:
			time.Sleep(duration)
		}

	}
}

type ErrorResponse struct {
	Error  string `json:"error" example:"error message"`
	Status string `json:"status" example:"error"`
}

func SendError(ctx context.Context, w http.ResponseWriter, errorMessage string, statusCode int) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Errorf("Отправлен код %d. ОШИБКА: %s", statusCode, errorMessage)

	response := ErrorResponse{Error: errorMessage, Status: "error"}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
