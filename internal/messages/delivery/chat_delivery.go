package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	jwt "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type MessageController struct {
	usecase usecase.MessageUsecase
}

func NewMessageController(usecase usecase.MessageUsecase) MessageController {
	return MessageController{
		usecase: usecase,
	}
}

func (h *MessageController) AddNewMessage(w http.ResponseWriter, r *http.Request) {
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting adding new message for chat: %v", chatUUID)

	if err != nil {
		//conn.400
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(jwt.User)
	log.Println(user)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	var messageDTO models.Message
	err = json.NewDecoder(r.Body).Decode(&messageDTO)

	if err != nil {
		log.Printf("Не удалось распарсить Json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.usecase.SendMessage(r.Context(), user.ID, chatUUID, messageDTO)

	if err != nil {
		log.Printf("Не удалось добавить сообщение: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MessageController) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId := mapVars["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	log.Println(mapVars["chatId"])
	log.Printf("Message Delivery: starting getting all messages for chat: %v", chatUUID)

	if err != nil {
		//conn.400
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := h.usecase.GetMessages(r.Context(), chatUUID, 0)
	if err != nil {
		log.Println("Error reading message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(messages)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	messageChannel := make(chan []models.Message, 10)
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
		case messages := <-messageChannel:
			// запись новых сообщений
			log.Println("Message delivery websocket: получены новые сообщения")

			if len(messages) > 0 {
				conn.WriteJSON(models.MessagesArrayDTOOutput{
					Messages: messages,
					IsNew:    true,
				})
			}
		default:
			time.Sleep(duration)
		}

	}
}
