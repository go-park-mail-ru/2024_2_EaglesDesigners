package controller

import (
	"encoding/json"
	"net/http"

	interfaces "github.com/go-park-mail-ru/2024_2_EaglesDesigner/login/interface"
)

var jwtSecret = []byte("КТо пРочитАл тОт сдОхНет :)")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

type LoginController struct {
	interfaces.ILoginService
}

func (controller *LoginController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusUnauthorized)
		return
	}

	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "invalid format JSON", http.StatusBadRequest)
		return
	}

	if controller.Authenticate(credentials.Username, credentials.Password) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Authentication successful")
	} else {
		http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
	}
}
