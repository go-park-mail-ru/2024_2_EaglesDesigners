package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/controller"
	chat "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list"
)

func main() {
	router := mux.NewRouter()

	// router.Use(cors.CorsMiddleware)
	router.MethodNotAllowedHandler = http.HandlerFunc(controller.MethodNotAllowedHandler)

	auth := auth.SetupController()
	chat := chat.SetupController()

	router.HandleFunc("/", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/auth", auth.AuthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/signup", auth.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/chats", auth.Middleware(chat.Handler)).Methods("GET", "OPTIONS")
	// http.HandleFunc("/logout")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:8001",
			"https://127.0.0.1:8001",
			"http://localhost:8001",
			"https://localhost:8001",
			"http://213.87.152.18:8001",
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
