package delivery

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/responser"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"
)

type Usecase interface {
	GetFile(ctx context.Context, fileIDStr string) (*bytes.Buffer, *models.FileMetaData, error)
	SaveSticker(ctx context.Context, file multipart.File, header *multipart.FileHeader, name string) error
	GetStickerPack(ctx context.Context, packID string) (models.GetStickerPackResponse, error)
	GetStickerPacks(ctx context.Context) (models.StickerPacks, error)
}

type Delivery struct {
	usecase Usecase
}

// /files/675f2ea013dbaf51a93aa2d3
// /files/675f466313dbaf51a93aa2e4
// /files/675f391413dbaf51a93aa2db.
func New(usecase Usecase) *Delivery {
	URLs := []string{
		"/uploads/stickers/675f2ea013dbaf51a93aa2d3.webp",
		"/uploads/stickers/675f466313dbaf51a93aa2e4.webp",
		"/uploads/stickers/675f391413dbaf51a93aa2db.webp",
		"/uploads/stickers/6762d25b5803e3d181d0ecc4.webp",
		"/uploads/stickers/6762d4535803e3d181d0ecc6.webp",
		"/uploads/stickers/6762d4545803e3d181d0ecc7.webp",
		"/uploads/stickers/6762d5135803e3d181d0ecc8.webp",

		"/uploads/stickers/6762d5505803e3d181d0ecc9.webp",
		"/uploads/stickers/6762d7f95803e3d181d0ecca.webp",
		"/uploads/stickers/6762d8aa5803e3d181d0eccb.webp",
		"/uploads/stickers/6762d8d85803e3d181d0eccc.webp",
		"/uploads/stickers/6762d8f45803e3d181d0eccd.webp",
		"/uploads/stickers/6762d90e5803e3d181d0ecce.webp",
		"/uploads/stickers/6762d9215803e3d181d0eccf.webp",
	}

	for _, url := range URLs {
		file, header, err := getMultipartFile(url)
		if err != nil {
			fmt.Errorf("stickers get error: %v", err)
		}

		err = usecase.SaveSticker(context.Background(), file, header, url)
		if err != nil {
			fmt.Errorf("stickers save error: %v", err)
		}
	}

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
// @Router /files/{fileID} [get].
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

// GetStickerPack godoc
// @Summary get sticker pack
// @Tags files
// @Param packid path string true "packid ID (UUID)"
// @Success 200 {object} models.GetStickerPackResponse "пак успешно получен"
// @Failure 404	{object} responser.ErrorResponse "Не найдено"
// @Router /stickerpacks/{packid} [get].
func (d Delivery) GetStickerPack(w http.ResponseWriter, r *http.Request) {
	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()

	vars := mux.Vars(r)
	packid := vars["packid"]

	pack, err := d.usecase.GetStickerPack(ctx, packid)
	if err != nil {
		log.WithError(err).Errorln("не удалось получить пак")
		responser.SendError(ctx, w, "not found", http.StatusNotFound)
		return
	}

	// responser.SendStruct(ctx, w, pack, http.StatusOK)
	jsonResp, err := easyjson.Marshal(pack)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusOK)
}

// GetStickerPacks godoc
// @Summary get sticker pack
// @Tags files
// @Success 200 {object} models.StickerPacks "паки успешно получены"
// @Failure 500	{object} responser.ErrorResponse "Внутреннее"
// @Router /stickerpacks [get].
func (d Delivery) GetStickerPacks(w http.ResponseWriter, r *http.Request) {
	log := logger.LoggerWithCtx(r.Context(), logger.Log)
	ctx := r.Context()

	packs, err := d.usecase.GetStickerPacks(ctx)
	if err != nil {
		log.WithError(err).Errorln("не удалось получить паки")
		responser.SendError(ctx, w, "internal server error", http.StatusInternalServerError)
		return
	}

	// responser.SendStruct(ctx, w, packs, http.StatusOK)
	jsonResp, err := easyjson.Marshal(packs)
	responser.SendJson(ctx, w, jsonResp, err, http.StatusOK)
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

type File struct {
	*os.File
}

// Реализация интерфейса multipart.File.
func (f *File) Close() error {
	return f.File.Close()
}

func newFileHeader(filePath string) (*multipart.FileHeader, error) {
	// Получаем информацию о файле
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Заполняем заголовок
	header := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	header.Header.Set("Content-Type", "image/webp")

	fmt.Print("sticker header: ", header)

	return header, nil
}

// Глобальная функция для получения файла и заголовка.
func getMultipartFile(filePath string) (multipart.File, *multipart.FileHeader, error) {
	// Открываем файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Errorf("файл не существует: %s", filePath)
		return nil, nil, fmt.Errorf("файл не существует: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("sticker error open: %v", err)
		return nil, nil, err
	}

	fmt.Println("handler sticker: ", file)

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		fmt.Errorf("нет инфы файла: %s", filePath)
		return nil, nil, err
	}
	if fileInfo.Size() == 0 {
		file.Close() // Закрываем файл при отсутствии данных
		fmt.Errorf("файл пустой: %s", filePath)
		return nil, nil, fmt.Errorf("файл пустой: %s", filePath)
	}

	// Получаем заголовок
	header, err := newFileHeader(filePath)
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		fmt.Errorf("sticker header error open: %v", err)
		return nil, nil, err
	}

	return &File{File: file}, header, nil
}
