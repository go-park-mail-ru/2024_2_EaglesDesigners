package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth/controller"
	chat "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/cors"
)

func main() {
	router := mux.NewRouter()

	router.Use(cors.CorsMiddleware)
	router.MethodNotAllowedHandler = http.HandlerFunc(controller.MethodNotAllowedHandler)

	auth := auth.SetupController()
	chat := chat.SetupController()

	router.HandleFunc("/", auth.AuthHandler).Methods("GET")
	router.HandleFunc("/auth", auth.AuthHandler).Methods("GET")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/chats", auth.Middleware(chat.Handler)).Methods("GET")
	// http.HandleFunc("/logout")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
