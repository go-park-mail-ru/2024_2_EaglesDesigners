package multiparthelper

import (
	"errors"
	"io"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	uploadPath = "/uploads/"
)

var ErrNotImage = errors.New("file is not image")

func SavePhoto(file multipart.File, folderName string) (string, error) {
	if ok := IsImageFile(file); !ok {
		return "", ErrNotImage
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
		return "", err
	}

	return path, nil
}

func RewritePhoto(file multipart.File, photoURL string) error {
	dst, err := os.Create(photoURL)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func RemovePhoto(photoURL string) error {
	err := os.Remove(photoURL)
	if err != nil {
		return err
	}

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
