package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/controller"
)

func main() {
	auth := auth.SetupController()

	http.HandleFunc("/", auth.AuthHandler)
	http.HandleFunc("/auth", auth.AuthHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/signup", auth.RegisterHandler)

	// http.HandleFunc("/logout")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
