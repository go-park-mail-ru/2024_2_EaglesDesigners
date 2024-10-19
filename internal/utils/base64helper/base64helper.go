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

	err = os.WriteFile("../../images/"+filenameUUID.String(), photoBytes, 0777)

	if err != nil {
		log.Printf("Unable to write into file %v: %v", filenameUUID, err)
		return uuid.Nil, err
	}
	return filenameUUID, nil
}

func ReadPhotoBase64(pgotoURL string) (string, error) {
	photoBytes, err := os.ReadFile("../../images/" + pgotoURL)
	if err != nil {
		log.Printf("Unable to read file %v: %v", pgotoURL, err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(photoBytes), nil
}
