package controller

import (
	"../service"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chats, err := service.GetChats(r.Cookies())
	if err != nil {
		fmt.Println(err)

		//вернуть 401

		return
	}

	fmt.Println(chats)
	//ебануть чаты в json и вернуть

	w.Write([]byte("{dwttq}"))
}

func init() {
	http.HandleFunc("/chats", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
