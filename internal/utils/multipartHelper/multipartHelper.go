package multiparthelper

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	uploadPath = "/uploads/"
)

// mb не нужон
func ReadPhoto(photoId uuid.UUID) ([]byte, error) {
	// file, err := os.Open("images/" + photoId.String() + ".png")
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	// data, err := io.ReadAll(file)
	// if err != nil {
	// 	return nil, err
	// }

	return []byte{}, nil
}

func SavePhoto(file multipart.File, folderName string) (string, error) {
	if ok := IsImageFile(file); !ok {
		return "", errors.New("file is not image")
	}

	filenameUUID := uuid.New()

	path := uploadPath + folderName + "/" + filenameUUID.String() + ".png"

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Unable to write into file %v: %v", filenameUUID, err)
		return "", err
	}

	log.Println("Фото сохранено")

	return path, nil
}

func RewritePhoto(file multipart.File, photoURL string) error {
	if ok := IsImageFile(file); !ok {
		return errors.New("file is not image")
	}

	log.Printf("Открытие файла %s\n", photoURL)

	dst, err := os.Create(photoURL)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Unable to rewrite into file %v: %v", photoURL, err)
		return err
	}

	log.Println("Фото перезаписано")
	return nil
}

func IsImageFile(file multipart.File) bool {
	img, err := imaging.Decode(file)
	if err != nil {
		return false
	}

	// сброс указателя на начало файла
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false
	}

	return img != nil
}
