package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	_ "github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/controller"
)

func main() {
	authController := auth.SetupController()

	// http.HandleFunc("/")
	http.HandleFunc("/login", authController.AuthHandler)
	http.HandleFunc("/register", authController.RegisterHandler)
	// http.HandleFunc("/logout")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
