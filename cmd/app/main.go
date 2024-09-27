package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	chat "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list"
)

func main() {
	router := mux.NewRouter()

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
