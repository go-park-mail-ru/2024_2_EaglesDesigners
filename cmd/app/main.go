package main

import (
	"log"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/docs"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	authDelivery "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/delivery"
	authRepo "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/repository"
	authUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/auth/usecase"
	chatController "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/controller"
	chatRepository "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/repository"
	chatService "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chat_list/service"
	tokenUC "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/responser"
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
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(responser.MethodNotAllowedHandler)

	authRepo := authRepo.NewRepository()
	tokenUC := tokenUC.NewUsecase(authRepo)
	authUC := authUC.NewUsecase(authRepo, tokenUC)
	auth := authDelivery.NewDelivery(authUC, tokenUC)

	chatRepo := chatRepository.NewChatRepository()
	chatService := chatService.NewChatService(tokenUC, chatRepo)
	chat := chatController.NewChatController(chatService)

	router.HandleFunc("/", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/auth", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/signup", auth.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/chats", auth.Middleware(chat.Handler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

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
		AllowedMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})

	handler := c.Handler(router)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
