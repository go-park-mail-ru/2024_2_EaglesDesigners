package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/docs"
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

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
	"github.com/redis/go-redis/v9"
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
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/patefon")
	// pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/patefon")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("База данных подключена")


	redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "1234",            
        DB:       0,              
    })
	status := redisClient.Ping(context.Background())

	if err := status.Err(); err != nil {
		log.Println("Не удалось подлключить Redis")
		return
	}
	defer redisClient.Close()
	log.Println("Redis подключен")

	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(responser.MethodNotAllowedHandler)

	// auth
	authRepo := authRepo.NewRepository(pool)
	tokenUC := tokenUC.NewUsecase(authRepo)
	authUC := authUC.NewUsecase(authRepo, tokenUC)
	auth := authDelivery.NewDelivery(authUC, tokenUC)

	// profile
	profileRepo := profileRepo.New(pool)
	profileUC := profileUC.New(profileRepo)
	profile := profileDelivery.New(profileUC, tokenUC)

	// chats
	messageRepo := messageRepository.NewMessageRepositoryImpl(pool)

	chatRepo, _ := chatRepository.NewChatRepository(pool)
	chatService := chatService.NewChatUsecase(tokenUC, chatRepo, messageRepo)
	chat := chatController.NewChatDelivery(chatService)

	// contacts
	contactsRepo := contactsRepo.New(pool)
	contactsUC := contactsUC.New(contactsRepo)
	contacts := contactsDelivery.New(contactsUC, tokenUC)

	// messages
	messageUsecase := messageUsecase.NewMessageUsecaseImpl(messageRepo, tokenUC, redisClient)
	messageDelivery := messageDelivery.NewMessageController(messageUsecase)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// r = r.WithContext(ctx) не работает nux.Vars(r), т.к. убирается сонтекст
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/auth", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/signup", auth.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/chats", auth.Middleware(chat.GetUserChatsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/addchat", auth.Middleware(chat.AddNewChat)).Methods("POST", "OPTIONS")
	router.HandleFunc("/addusers", auth.Middleware(chat.AddUsersIntoChat)).Methods("POST", "OPTIONS")

	router.HandleFunc("/addusers", auth.Middleware(chat.AddUsersIntoChat)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", auth.Middleware(profile.GetProfileHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile", auth.Middleware(profile.UpdateProfileHandler)).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/chats", auth.Middleware(chat.Handler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/contacts", auth.Middleware(contacts.GetContactsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/contacts", auth.Middleware(contacts.AddContactHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	tmpl := template.Must(template.ParseFiles("index.html"))

	router.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	router.HandleFunc("/chat/{chatId}/messages", auth.Middleware(messageDelivery.GetAllMessages)).Methods("GET", "OPTIONS")
	router.HandleFunc("/chat/{chatId}/messages", auth.Middleware(messageDelivery.AddNewMessage)).Methods("POST", "OPTIONS")
	router.HandleFunc("/chat/{chatId}", auth.Middleware(messageDelivery.HandleConnection))

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
