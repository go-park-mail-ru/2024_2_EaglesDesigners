package delivery

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/middleware"
	websocketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/usecase"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 10048,

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
	// 		"https://patefon.site",
	// 		"http://localhost",
	// 		"https://localhost",
	// 		"https://localhost:8083",
	// 		"http://localhost:8083",
	// 		"http://localhost:9090",
	// 		"https://localhost:9090",
	// 		"http://127.0.0.1:9090",
	// 		"https://127.0.0.1:9090",
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
	metric.IncHit()
	log := logger.LoggerWithCtx(r.Context(), logger.Log)

	id, _ := uuid.Parse("fa4e08e4-1024-49cb-a799-4aa2a4f3a9df")

	user := middleware.User{
		ID: id,
	}
	log.Printf("Пользователь %v Открыл сокет", user.ID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		responser.SendError(r.Context(), w, "Delivery: error during connection upgrade:", http.StatusInternalServerError)
		return
	}
	defer log.Println("Message delivery: websocket is closing")
	defer conn.Close()

	eventChannel := make(chan websocketUsecase.AnyEvent, 10)

	err = h.usecase.InitBrokersForUser(user.ID, eventChannel)
	if err != nil {
		log.Errorf("Не удалось иницировать брокеры для пользователя")
		responser.SendError(r.Context(), w, "Нет нужных параметров", http.StatusInternalServerError)

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
