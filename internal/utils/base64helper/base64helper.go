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

	err = os.WriteFile("../../../images/"+filenameUUID.String() + ".png", photoBytes, 0777)

	if err != nil {
		log.Printf("Unable to write into file %v: %v", filenameUUID, err)
		return uuid.Nil, err
	}
	return filenameUUID, nil
}


func ReadPhotoBase64(photoId uuid.UUID) (string, error) {
	photoBytes, err := os.ReadFile("images/" + photoId.String() + ".png")
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

