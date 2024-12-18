package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/authv1"
	authDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/middleware"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/delivery"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/websocket_service/internal/websocket/usecase"
)

const (
	host = "patefon"
	port = 8082
)

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

	socketUsecase := usecase.NewWebsocketUsecase(ch, host, port)
	socketDelivery := delivery.NewWebsocket(*socketUsecase)

	router := mux.NewRouter()

	router = router.PathPrefix("/api/").Subrouter()

	// auth

	grpcConnAuth, err := grpc.Dial(
		"auth:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConnAuth.Close()
	authClient := authv1.NewAuthClient(grpcConnAuth)

	auth := authDelivery.New(authClient)

	// ручки
	tmpl := template.Must(template.ParseFiles("index.html"))

	router.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	router.HandleFunc("/startwebsocket", auth.Authorize(socketDelivery.HandleConnection))
	// мктрики
	router.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8083")
	if err := http.ListenAndServe(":8083", router); err != nil {
		log.Fatal(err)
	}
}
