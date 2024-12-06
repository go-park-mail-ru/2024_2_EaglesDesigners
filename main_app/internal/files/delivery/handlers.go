package delivery

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"
	"github.com/gorilla/mux"
)

type Usecase interface {
	GetFile(ctx context.Context, fileIDStr string) (*bytes.Buffer, *models.FileMetaData, error)
}

type Delivery struct {
	usecase Usecase
}

func New(usecase Usecase) *Delivery {
	return &Delivery{
		usecase: usecase,
	}
}

// // GetImage godoc
// // @Summary Retrieve an image
// // @Description Fetches an image from the specified folder and by filename
// // @Tags uploads
// // @Accept json
// // @Produce json
// // @Param folder path string true "Folder name" example("avatar")
// // @Param name path string true "File name" example("642c5a57-ebc7-49d0-ac2d-f2f1f474bee7.png")
// // @Success 200 {file} string "Successful image retrieval"
// // @Failure 404 {object} map[string]string "File not found"
// // @Router /files/images/{name} [get]
// func (d *Delivery) GetImage(w http.ResponseWriter, r *http.Request) {
// 	log := logger.LoggerWithCtx(r.Context(), logger.Log)

// 	vars := mux.Vars(r)
// 	folder := vars["folder"]
// 	name := vars["name"]

// 	imagePath := "/uploads/" + folder + "/" + name

// 	log.Println("пришел запрос на получение картинки ", imagePath)

// 	http.ServeFile(w, r, imagePath)
// }

// // GetImage godoc
// // @Summary Получить файл
// // @Description Получить файл по его Id
// // @Tags files
// // @Accept json
// // @Produce json
// // @Param fileID path string true "File ID" example("642c5a57-ebc7-49d0-ac2d-f2f1f474bee7")
// // @Success 200 {file} responser.SuccessResponse "Файл успешно получен"
// // @Failure 404 {object} responser.ErrorResponse "файл не найден"
// @Router /files/{fileID} [get]
func (d Delivery) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	vars := mux.Vars(r)
	fileIDStr := vars["fileID"]

	log.Printf("fileID: %s", fileIDStr)

	fileBuffer, metadata, err := d.usecase.GetFile(ctx, fileIDStr)
	if err != nil {
		log.WithError(err).Errorln("не удалось получить файл")
		responser.SendError(ctx, w, "File not found", http.StatusNotFound)
		return
	}

	if metadata != nil {
		w.Header().Set("Content-Type", metadata.ContentType)
		w.Header().Set("Content-Disposition", "attachment; filename=\""+metadata.Filename+"\"")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", metadata.FileSize))
		w.WriteHeader(http.StatusOK)
	} else {
		log.Warnln("нет метаданных")
	}

	if _, err := w.Write(fileBuffer.Bytes()); err != nil {
		http.Error(w, "Could not send file", http.StatusInternalServerError)
		return
	}
}

// // @Router /files [post]
// func (d *Delivery) UploadFile(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	filename, err := d.usecase.SaveFile(ctx, file, header)
// 	if err != nil {
// 		responser.SendError(ctx, w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	responser.SendStruct(ctx, w, models.UploadFileResponse{FileID: filename}, http.StatusCreated)
// }
