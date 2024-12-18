package delivery

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
)

type Delivery struct{}

func New() *Delivery {
	return &Delivery{}
}

// GetImage godoc
// @Summary Retrieve an image
// @Description Fetches an image from the specified folder and by filename
// @Tags uploads
// @Accept json
// @Produce json
// @Param folder path string true "Folder name" example("avatar")
// @Param name path string true "File name" example("642c5a57-ebc7-49d0-ac2d-f2f1f474bee7.png")
// @Success 200 {file} string "Successful image retrieval"
// @Failure 404 {object} map[string]string "File not found"
// @Router /uploads/{folder}/{name} [get].
func (d *Delivery) GetImage(w http.ResponseWriter, r *http.Request) {
	log := logger.LoggerWithCtx(r.Context(), logger.Log)

	vars := mux.Vars(r)
	folder := vars["folder"]
	name := vars["name"]

	imagePath := "/uploads/" + folder + "/" + name

	log.Println("пришел запрос на получение картинки ", imagePath)

	http.ServeFile(w, r, imagePath)
}
