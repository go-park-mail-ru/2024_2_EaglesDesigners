package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"strings"

	"github.com/chai2010/webp"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const fileURLPrefix = "/files/"

type Repository interface {
	GetFile(ctx context.Context, filename string) (*bytes.Buffer, *models.FileMetaData, error)
	SaveFile(ctx context.Context, fileBuffer *bytes.Buffer, metadata primitive.D) (string, error)
	DeleteFile(ctx context.Context, fileID primitive.ObjectID) error
	RewriteFile(ctx context.Context, fileID primitive.ObjectID, fileBuffer *bytes.Buffer, metadata primitive.D) error
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) GetFile(ctx context.Context, fileIDStr string) (*bytes.Buffer, *models.FileMetaData, error) {
	return u.repo.GetFile(ctx, fileIDStr)
}

func (u *Usecase) SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (models.Payload, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("сохранение файла для пользователей: ", users)

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать буфер")
		return models.Payload{}, err
	}

	metadata := getFileMetadata(header)

	if len(users) > 0 {
		metadata = append(metadata, primitive.E{Key: "users", Value: users})
	}

	fileID, err := u.repo.SaveFile(ctx, &fileBuffer, metadata)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать файл")
		return models.Payload{}, err
	}

	out := getFileNameAndSize(header)
	out.URL = addFileURLPrefix(fileID)

	return out, nil
}

func (u *Usecase) SavePhoto(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (models.Payload, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	contentType, err := isImage(*header)
	if err != nil {
		log.WithError(err).Errorln("файл не фото")
		return models.Payload{}, err
	}

	fileBuffer, err := convertToWebP(file, contentType)
	if err != nil {
		log.WithError(err).Errorln("не удалось конвертировать фото в webp")
		return models.Payload{}, err
	}

	metadata := getPhotoMetadata(header, int64(fileBuffer.Len()))

	if len(users) > 0 {
		metadata = append(metadata, primitive.E{Key: "users", Value: users})
	}

	fileID, err := u.repo.SaveFile(ctx, &fileBuffer, metadata)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать файл")
		return models.Payload{}, err
	}

	out := getFileNameAndSize(header)
	out.URL = addFileURLPrefix(fileID)

	return out, nil
}

// for chats
func (u *Usecase) SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	contentType, err := isImage(*header)
	if err != nil {
		log.WithError(err).Errorln("файл не фото")
		return "", err
	}

	fileBuffer, err := convertToWebP(file, contentType)
	if err != nil {
		log.WithError(err).Errorln("не удалось конвертировать фото в webp")
		return "", err
	}

	metadata := getPhotoMetadata(header, int64(fileBuffer.Len()))

	fileID, err := u.repo.SaveFile(ctx, &fileBuffer, metadata)

	return addFileURLPrefix(fileID), err
}

// for profile
func (u *Usecase) RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Printf("пришел запрос на перезапись %s", fileIDStr)

	parts := strings.Split(fileIDStr, "/")
	var ID string
	if len(parts) > 2 {
		ID = parts[2]
	} else {
		log.Errorln("ID не найден")
		return errors.New("ID не найден")
	}

	contentType, err := isImage(header)
	if err != nil {
		log.WithError(err).Errorln("файл не фото")
		return err
	}

	fileID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.WithError(err).Errorln("не удалось преобразовать в objectID")
		return err
	}

	fileBuffer, err := convertToWebP(file, contentType)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать буфер")
		return err
	}

	metadata := getPhotoMetadata(&header, int64(fileBuffer.Len()))

	err = u.repo.RewriteFile(ctx, fileID, &fileBuffer, metadata)
	if err != nil {
		log.WithError(err).Errorln("не удалось перезаписать файл")
		return err
	}

	return nil
}

func (u *Usecase) DeletePhoto(ctx context.Context, fileIDStr string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Printf("пришел запрос на удаление %s", fileIDStr)

	parts := strings.Split(fileIDStr, "/")
	var ID string
	if len(parts) > 2 {
		ID = parts[2]
	} else {
		log.Errorln("ID не найден")
		return errors.New("ID не найден")
	}

	fileID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.WithError(err).Errorln("не удалось преобразовать в objectID")
		return err
	}

	return u.repo.DeleteFile(ctx, fileID)
}

// func (u *Usecase) RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) (string, error) {
// 	log := logger.LoggerWithCtx(ctx, logger.Log)

// 	err := isImage(header)
// 	if err != nil {
// 		log.WithError(err).Errorln("файл не фото")
// 		return "", err
// 	}

// 	fileID, err := primitive.ObjectIDFromHex(fileIDStr)
// 	if err != nil {
// 		log.WithError(err).Errorln("не удалось преобразовать в objectID")
// 		return "", err
// 	}

// 	fileBuffer, err := getFileBuffer(file)
// 	if err != nil {
// 		log.WithError(err).Errorln("не удалось создать буфер")
// 		return "", err
// 	}

// 	metadata := getFileMetadata(&header)

// 	newFileID, err := u.repo.SaveFile(ctx, &fileBuffer, metadata)
// 	if err != nil {
// 		log.WithError(err).Errorln("не удалось создать файл")
// 		return "", err
// 	}

// 	err = u.repo.DeleteFile(ctx, fileID)
// 	if err != nil {
// 		log.WithError(err).Warnln("не удалось удалить файл")
// 	}

// 	return newFileID, nil
// }

func (u *Usecase) UpdateFile(ctx context.Context, fileIDStr string, file multipart.File, header *multipart.FileHeader) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	fileID, err := primitive.ObjectIDFromHex(fileIDStr)
	if err != nil {
		log.WithError(err).Errorln("не удалось преобразовать в objectID")
		return "", err
	}

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать буфер")
		return "", err
	}

	metadata := getFileMetadata(header)

	newFileID, err := u.repo.SaveFile(ctx, &fileBuffer, metadata)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать файл")
		return "", err
	}

	err = u.repo.DeleteFile(ctx, fileID)
	if err != nil {
		log.WithError(err).Warnln("не удалось удалить файл")
	}

	return newFileID, nil
}

func getFileNameAndSize(header *multipart.FileHeader) models.Payload {
	return models.Payload{
		Filename: header.Filename,
		Size:     header.Size,
	}
}

func getFileMetadata(header *multipart.FileHeader) bson.D {
	return bson.D{
		{Key: "filename", Value: header.Filename},
		{Key: "contentType", Value: header.Header.Get("Content-Type")},
		{Key: "size", Value: header.Size},
	}
}

func getPhotoMetadata(header *multipart.FileHeader, size int64) bson.D {
	return bson.D{
		{Key: "filename", Value: header.Filename},
		{Key: "contentType", Value: "image/webp"},
		{Key: "size", Value: size},
	}
}

func getFileBuffer(file multipart.File) (bytes.Buffer, error) {
	fileBuffer := new(bytes.Buffer)

	if _, err := io.Copy(fileBuffer, file); err != nil {
		return bytes.Buffer{}, err
	}

	return *fileBuffer, nil
}

func (u *Usecase) IsImage(header multipart.FileHeader) error {
	_, err := isImage(header)
	return err
}

func isImage(header multipart.FileHeader) (string, error) {
	imageType := header.Header.Get("Content-Type")
	switch imageType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return imageType, nil
	default:
		return "", fmt.Errorf("недопустимый тип файла: %s", header.Header.Get("Content-Type"))
	}
}

func addFileURLPrefix(fileID string) string {
	return fileURLPrefix + fileID
}

func convertToWebP(file multipart.File, contentType string) (bytes.Buffer, error) {
	var img image.Image
	var err error

	switch contentType {
	case "image/gif":
		img, err = gif.Decode(file)
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	case "image/webp":
		img, err = webp.Decode(file)
	default:
		return bytes.Buffer{}, fmt.Errorf("unsupported image format: %s", contentType)
	}

	if err != nil {
		return bytes.Buffer{}, err
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 70}); err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}
