package delivery

import (
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	websocketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/websocket/usecase"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5024,
	WriteBufferSize: 5024,

	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := []string{
			"http://127.0.0.1:8001",
			"https://127.0.0.1:8001",
			"http://localhost:8001",
			"https://localhost:8001",
			"http://213.87.152.18:8001",
			"http://212.233.98.59:8001",
			"https://213.87.152.18:8001",
			"http://212.233.98.59:8080",
			"https://212.233.98.59:8080",
		}

		for _, origin := range allowedOrigins {
			if r.Header.Get("Origin") == origin {
				return true
			}
		}
		return false
	},
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

	user, ok := r.Context().Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		log.Println("Message delivery -> AddNewMessage: нет юзера в контексте")
		responser.SendError(r.Context(), w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
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

			conn.WriteJSON(message)

		default:
			time.Sleep(duration)
		}

	}
}
