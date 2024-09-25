package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/auth"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/service"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/src/chat_list/models"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Пришёл запрос на получения чатов")

	chats, err := service.GetChats(r.Cookies())
	if err != nil {
		fmt.Println(err)

		//вернуть 401
		w.WriteHeader(http.StatusUnauthorized)

		log.Printf("НЕ УДАЛОСЬ ПОЛУЧИТЬ ЧАТЫ. ОШИБКА: %s", err)
		return
	}
	jsonResp, err := json.Marshal(chats)


	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func init() {
	auth := auth.SetupController()
	http.HandleFunc("/chats", auth.Middleware(handler))
	// http.HandleFunc("/chats", handler)
}
