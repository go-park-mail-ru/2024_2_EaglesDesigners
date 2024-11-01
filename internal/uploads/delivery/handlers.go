package delivery

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Delivery struct {
}

func New() *Delivery {
	return &Delivery{}
}

func (d *Delivery) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	name := vars["name"]

	imagePath := "/uploads/" + folder + "/" + name

	log.Println("Uploads delivery: пришел запрос на получение картинки ", imagePath)

	http.ServeFile(w, r, imagePath)
}
