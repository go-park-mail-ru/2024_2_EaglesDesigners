package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/docs"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	authDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/delivery"
	authRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	authUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	chatController "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/delivery"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/repository"
	chatService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/usecase"
	contactsDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/delivery"
	contactsRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/repository"
	contactsUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/usecase"
	tokenUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	messageDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/delivery"
	messageRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/repository"
	messageUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/usecase"
	profileDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/delivery"
	profileRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/repository"
	profileUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/usecase"
	uploadsDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/uploads/delivery"

	websocketDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/websocket/delivery"
	websocketUsecase "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/websocket/usecase"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
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
	authRepo := authRepo.NewRepository(pool)
	tokenUC := tokenUC.NewUsecase(authRepo)
	authUC := authUC.NewUsecase(authRepo, tokenUC)
	auth := authDelivery.NewDelivery(authUC, tokenUC)

	// uploads

	uploads := uploadsDelivery.New()

	// profile
	profileRepo := profileRepo.New(pool)
	profileUC := profileUC.New(profileRepo)
	profile := profileDelivery.New(profileUC, tokenUC)

	// chats
	messageRepo := messageRepository.NewMessageRepositoryImpl(pool)

	chatRepo, _ := chatRepository.NewChatRepository(pool)

	messageUsecase := messageUsecase.NewMessageUsecaseImpl(messageRepo, chatRepo, tokenUC, ch)

	chatService := chatService.NewChatUsecase(tokenUC, chatRepo, messageRepo, ch)
	chat := chatController.NewChatDelivery(chatService)

	// contacts
	contactsRepo := contactsRepo.New(pool)
	contactsUC := contactsUC.New(contactsRepo)
	contacts := contactsDelivery.New(contactsUC, tokenUC)

	// messages

	messageDelivery := messageDelivery.NewMessageController(messageUsecase)

	// websocket
	websocketUsecase := websocketUsecase.NewWebsocketUsecase(ch, chatRepo)
	websocketDelivery := websocketDelivery.NewWebsocket(*websocketUsecase)

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

	tmpl := template.Must(template.ParseFiles("index.html"))

	router.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

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

	router.HandleFunc("/chat/{chatId}/messages/search", auth.Authorize(auth.Csrf(messageDelivery.SearchMessages))).Methods("GET", "OPTIONS")

	router.HandleFunc("/messages/{messageId}", auth.Authorize(auth.Csrf(messageDelivery.DeleteMessage))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/messages/{messageId}", auth.Authorize(auth.Csrf(messageDelivery.UpdateMessage))).Methods("PUT", "OPTIONS")

	router.HandleFunc("/startwebsocket", auth.Authorize(websocketDelivery.HandleConnection))

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
			"https://212.233.98.59:8080"},
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
