package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/delivery"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/usecase"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

const host = "localhost"
const port = 8090

func main() {
	// подключаем rebbit mq
	conn, err := amqp.Dial("amqp://root:root@rabbitmq:5672/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}
	defer func() {
		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
	}()
	log.Println("rebbit mq подключен")


	socketUsecase := usecase.NewWebsocketUsecase(ch, "localhost", port)

	socketDelivery := delivery.NewWebsocket(socketUsecase)

	router := mux.NewRouter()
	
	router.HandleFunc("/startwebsocket", auth.Authorize(socketDelivery.HandleConnection))
}
