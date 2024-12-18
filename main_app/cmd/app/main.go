package main

import (
	"context"
	"log"
	"net"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/docs"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/delivery"
	chatController "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/delivery"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/repository"
	chatService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/usecase"
	contactsDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/delivery"
	contactsRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/repository"
	contactsUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/contacts/usecase"
	messageDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/delivery"
	messageRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/repository"
	messageUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/usecase"
	profileDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/delivery"
	profileRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/repository"
	profileUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/usecase"
	uploadsDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/uploads/delivery"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/authv1"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"

	"github.com/asaskevich/govalidator"
	amqp "github.com/rabbitmq/amqp091-go"
)

// swag init

// @title           Swagger Patefon API
// @version         1.0
// @description     This is a description of the Patefon server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      212.233.98.59:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@postgres:5432/patefon")
	// pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/patefon")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("База данных подключена")

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

	govalidator.SetFieldsRequiredByDefault(true)

	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(responser.MethodNotAllowedHandler)

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

	// token рудемент

	// uploads

	uploads := uploadsDelivery.New()

	// profile
	profileRepo := profileRepo.New(pool)
	profileUC := profileUC.New(profileRepo)
	profile := profileDelivery.New(profileUC)

	// chats
	messageRepo := messageRepository.NewMessageRepositoryImpl(pool)

	chatRepo, _ := chatRepository.NewChatRepository(pool)

	messageUsecase := messageUsecase.NewMessageUsecaseImpl(messageRepo, chatRepo, ch)

	chatService := chatService.NewChatUsecase(chatRepo, messageRepo, ch)
	chat := chatController.NewChatDelivery(chatService)

	// contacts
	contactsRepo := contactsRepo.New(pool)
	contactsUC := contactsUC.New(contactsRepo)
	contacts := contactsDelivery.New(contactsUC)

	// messages

	messageDelivery := messageDelivery.NewMessageController(messageUsecase)

	// добавление request_id в ctx всем запросам
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			ctx = context.WithValue(r.Context(), logger.RequestIDKey, requestID)
			r = r.WithContext(ctx)

			log := logger.LoggerWithCtx(ctx, logger.Log)

			log.Printf("Пришел запрос %s", r.URL.String())

			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/", auth.Authorize(auth.AuthHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/auth", auth.Authorize(auth.AuthHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/signup", auth.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/uploads/{folder}/{name}", uploads.GetImage).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", auth.Authorize(profile.GetSelfProfileHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", auth.Authorize(auth.Csrf(profile.UpdateProfileHandler))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/profile/{userid}", profile.GetProfileHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/contacts", auth.Authorize(contacts.GetContactsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/contacts", auth.Authorize(auth.Csrf(contacts.AddContactHandler))).Methods("POST", "OPTIONS")
	router.HandleFunc("/contacts", auth.Authorize(auth.Csrf(contacts.DeleteContactHandler))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/contacts/search", auth.Authorize(contacts.SearchContactsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/chats", auth.Authorize(chat.GetUserChatsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/addchat", auth.Authorize(auth.Csrf(chat.AddNewChat))).Methods("POST", "OPTIONS")
	router.HandleFunc("/chat/search", auth.Authorize(auth.Csrf(chat.SearchChats))).Methods("GET", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/addusers", auth.Authorize(auth.Csrf(chat.AddUsersIntoChat))).Methods("POST", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/delusers", auth.Authorize(auth.Csrf(chat.DeleteUsersFromChat))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/deluser/{userId}", auth.Authorize(auth.Csrf(chat.DeleteUserFromChat))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/delete", auth.Authorize(auth.Csrf(chat.DeleteChatOrGroup))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/chat/{chatId}", auth.Authorize(auth.Csrf(chat.UpdateGroup))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/chat/{chatId}", auth.Authorize(chat.GetChatInfo)).Methods("GET", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/messages", auth.Authorize(messageDelivery.GetAllMessages)).Methods("GET", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/messages/pages/{lastMessageId}", auth.Authorize(messageDelivery.GetMessagesWithPage)).Methods("GET", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/messages", auth.Authorize(auth.Csrf(messageDelivery.AddNewMessage))).Methods("POST", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/{messageId}/branch", auth.Authorize(auth.Csrf(chat.AddBranch))).Methods("POST", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/leave", auth.Authorize(auth.Csrf(chat.LeaveChat))).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/channel/{channelId}/join", auth.Authorize(chat.JoinChannel)).Methods("POST", "OPTIONS")

	router.HandleFunc("/chat/{chatId}/messages/search", auth.Authorize(auth.Csrf(messageDelivery.SearchMessages))).Methods("GET", "OPTIONS")

	router.HandleFunc("/messages/{messageId}", auth.Authorize(auth.Csrf(messageDelivery.DeleteMessage))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/messages/{messageId}", auth.Authorize(auth.Csrf(messageDelivery.UpdateMessage))).Methods("PUT", "OPTIONS")

	// мктрики
	router.Handle("/metrics", promhttp.Handler())
	metric.RecordMetrics()

	// хз чо это
	http.HandleFunc("/", httpSwagger.Handler())

	go startMainServer(router)
	go startChatServerGRPC(chatService)

	select {}
}

func startMainServer(router *mux.Router) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:8001",
			"https://127.0.0.1:8001",
			"http://localhost:8001",
			"https://localhost:8001",
			"http://213.87.152.18:8001",
			"http://212.233.98.59:8001",
			"https://213.87.152.18:8001",
			"http://212.233.98.59:8080",
			"https://212.233.98.59:8080",
			"https://localhost:8083",
			"http://localhost:8083",
			"http://localhost:9090",
			"https://localhost:9090",
			"http://127.0.0.1:9090",
			"https://127.0.0.1:9090",
		},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})

	handler := c.Handler(router)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}

}

func startChatServerGRPC(chatService chatService.ChatUsecase) {
	// grpc for chat
	chatServer := grpc.NewServer()
	chatController.Register(chatServer, chatService)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting chat server at :8082")
	if err := chatServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
	log.Println("server started at :8082")
}
