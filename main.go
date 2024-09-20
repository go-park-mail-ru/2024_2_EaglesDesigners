package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/login"
)

func main() {
	loginController := login.SetupController()

	http.HandleFunc("/login", loginController.LoginHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
