package delivery

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (h *MessageController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	chatId := mux.Vars(r)["chatId"]
	chatUUID, err := uuid.Parse(chatId)

	if err != nil {
		//conn.400
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Delivery: error during connection upgrade:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Здесь можно хранить список старых сообщений (например, в массиве или в базе данных)
	messageChannel := make(chan []models.Message, 10)
	errChannel := make(chan error, 10)
	closeChannel := make(chan bool, 1)

	defer func() {
		closeChannel <- true
		close(closeChannel)
	}()

	go h.usecase.ScanForNewMessages(messageChannel, chatUUID, errChannel, closeChannel)

	// история чата
	messages, err := h.usecase.GetMessages(r.Context(), chatUUID, 0)

	if err != nil {
		log.Println("Error reading message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn.WriteJSON(models.MessagesArrayDTOOutput{
		Messages: messages.Messages,
		IsNew:    false,
	})

	// пока соеденено
	for {
		var msg models.MessageDTOInput
		err := conn.ReadJSON(&msg)

		// сообщение кривого формата
		if err != nil {
			log.Println("Error reading message:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// отправили запрос на дисконект
		if msg.Disconnect {
			log.Println("Delivery: close connection 200")
			w.WriteHeader(http.StatusOK)
			return
		}

		if msg.Message != "" {
			// если есть сообщение

			message := models.Message{
				Message: msg.Message,
			}

			err = h.usecase.SendMessage(r.Context(), chatUUID, message)
			if err != nil {
				log.Printf("Delivery: не удолось отправить сообщение: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("Delivery: сообщение успешно добавлено")
		}

		err = <-errChannel
		if err != nil {
			log.Printf("Delivery: ошибка в поиске новых сообщений: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages := <-messageChannel

		if len(messages) > 0 {
			conn.WriteJSON(models.MessagesArrayDTOOutput{
				Messages: messages,
				IsNew:    true,
			})
		}
	}
}
