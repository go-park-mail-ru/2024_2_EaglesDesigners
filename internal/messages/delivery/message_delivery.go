package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

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

type MessageController struct {
	usecase usecase.MessageUsecase
}

func NewMessageController(usecase usecase.MessageUsecase) MessageController {
	return MessageController{
		usecase: usecase,
	}
}

// AddNewMessageHandler godoc
// @Summary Add new message
// @Tags message
// @Accept json
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param message body models.MessageInput true "Message info"
// @Success 201 "Сообщение успешно добавлено"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось добавить сообщение"
// @Router /chat/{chatId}/messages [post]
func (h *MessageController) AddNewMessage(w http.ResponseWriter, r *http.Request) {
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	if err != nil {
		//conn.400
		log.Println("Delivery: error during parsing json:", err)
		responser.SendError(w, fmt.Sprintf("Delivery: error during connection upgrade:%v", err), http.StatusBadRequest)
		return
	}

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting adding new message for chat: %v", chatUUID)

	user, ok := r.Context().Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		log.Println("Message delivery -> AddNewMessage: нет юзера в контексте")
		responser.SendError(w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	var messageDTO models.Message
	err = json.NewDecoder(r.Body).Decode(&messageDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		responser.SendError(w, fmt.Sprintf("Не удалось распарсить Json: %v", err), http.StatusBadRequest)
		return
	}

	err = h.usecase.SendMessage(r.Context(), user, chatUUID, messageDTO)

	if err != nil {
		log.Printf("Не удалось добавить сообщение: %v", err)
		responser.SendError(w, fmt.Sprintf("Не удалось добавить сообщение: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllMessages godoc
// @Summary Add new message
// @Tags message
// @Param chatId path string true "Chat ID (UUID)" minlength(36) maxlength(36) example("123e4567-e89b-12d3-a456-426614174000")
// @Param message body models.MessagesArrayDTO true "Messages"
// @Success 200 "Сообщение успешно отаправлены"
// @Failure 400	{object} responser.ErrorResponse "Некорректный запрос"
// @Failure 500	{object} responser.ErrorResponse "Не удалось получить сообщениея"
// @Router /chat/{chatId}/messages [get]
func (h *MessageController) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	if err != nil {
		//conn.400
		log.Printf("Message delivery -> GetAllMessages: получен кривой Id юзера: %v", err)
		responser.SendError(w, fmt.Sprintf("Delivery: error during connection upgrade:%v", err), http.StatusBadRequest)
		return
	}

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting getting all messages for chat: %v", chatUUID)

	messages, err := h.usecase.GetMessages(r.Context(), chatUUID)
	if err != nil {
		log.Println("Error reading message:", err)
		responser.SendError(w, fmt.Sprintf("Error reading message:%v", err), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(messages)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		responser.SendError(w, fmt.Sprintf("error happened in JSON marshal. Err: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (h *MessageController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// начало

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer log.Println("Message delivery: websocket is closing")
	defer conn.Close()

	// Здесь можно хранить список старых сообщений (например, в массиве или в базе данных)
	messageChannel := make(chan models.WebScoketDTO, 10)
	errChannel := make(chan error, 10)
	closeChannel := make(chan bool, 1)

	defer func() {
		closeChannel <- true
		close(closeChannel)
	}()

	if err != nil {
		log.Println("Error reading message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go h.usecase.ScanForNewMessages(r.Context(), messageChannel, errChannel, closeChannel)

	// пока соеденено
	duration := 500 * time.Millisecond

	for {
		select {
		case err = <-errChannel:

			if err != nil {
				log.Printf("Delivery: ошибка в поиске новых сообщений: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case message := <-messageChannel:
			// запись новых сообщений
			log.Println("Message delivery websocket: получены новые сообщения")

			conn.WriteJSON(message)

		default:
			time.Sleep(duration)
		}

	}
}

