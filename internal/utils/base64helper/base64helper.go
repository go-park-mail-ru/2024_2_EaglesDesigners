package base64helper

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/google/uuid"
)

func SavePhotoBase64(base64Photo string) (uuid.UUID, error) {

	photoBytes, err := base64.StdEncoding.DecodeString(base64Photo)
	if err != nil {
		log.Printf("Ну удалось расшифровать фото: %v \n", err)
		return uuid.Nil, err
	}

	filenameUUID := uuid.New()

	//здесь можно fullpath. Если указывать relative, то от диркутории в которой пишете go run
	err = os.WriteFile("images/"+filenameUUID.String()+".png", photoBytes, 0777)

	if err != nil {
		log.Printf("Unable to write into file %v: %v", filenameUUID, err)
		return uuid.Nil, err
	}

	log.Println("Фото сохранено")
	return filenameUUID, nil
}

func ReadPhotoBase64(pgotoId uuid.UUID) (string, error) {
	photoBytes, err := os.ReadFile("images/" + pgotoId.String())

	if err != nil {
		log.Printf("Unable to read file %v: %v", photoId, err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(photoBytes), nil
}

// RewritePhoto перезапишет base64 фото в images по названию файла filename.
func RewritePhoto(base64Photo string, filename string) error {
	photoBytes, err := base64.StdEncoding.DecodeString(base64Photo)
	if err != nil {
		log.Printf("Ну удалось расшифровать фото: %v \n", err)
		return err
	}

	err = os.WriteFile("../../images/"+filename, photoBytes, 0777)

	if err != nil {
		log.Printf("Не удалось перезаписать фото %v: %v", filename, err)
		return err
	}
	return nil
}

