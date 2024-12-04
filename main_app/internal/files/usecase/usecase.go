package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	GetFile(ctx context.Context, fileID primitive.ObjectID) (*bytes.Buffer, *models.FileMetaData, error)
	SaveFile(ctx context.Context, fileBuffer *bytes.Buffer, metadata primitive.D) (string, error)
	DeleteFile(ctx context.Context, fileID primitive.ObjectID) error
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
	log := logger.LoggerWithCtx(ctx, logger.Log)

	fileID, err := primitive.ObjectIDFromHex(fileIDStr)
	if err != nil {
		log.WithError(err).Errorln("не удалось преобразовать в objectID")
		return nil, nil, err
	}
	return u.repo.GetFile(ctx, fileID)
}

func (u *Usecase) SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать буфер")
		return "", err
	}

	metadata := getFileMetadata(header)

	return u.repo.SaveFile(ctx, &fileBuffer, metadata)
}

func (u *Usecase) SavePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	err := isImage(header)
	if err != nil {
		log.WithError(err).Errorln("файл не фото")
		return "", err
	}

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		log.WithError(err).Errorln("не удалось создать буфер")
		return "", err
	}

	metadata := getFileMetadata(&header)

	return u.repo.SaveFile(ctx, &fileBuffer, metadata)
}

func (u *Usecase) RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	err := isImage(header)
	if err != nil {
		log.WithError(err).Errorln("файл не фото")
		return "", err
	}

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

	metadata := getFileMetadata(&header)

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

func getFileMetadata(header *multipart.FileHeader) bson.D {
	return bson.D{
		{Key: "filename", Value: header.Filename},
		{Key: "contentType", Value: header.Header.Get("Content-Type")},
		{Key: "size", Value: header.Size},
	}
}

func getFileBuffer(file multipart.File) (bytes.Buffer, error) {
	fileBuffer := new(bytes.Buffer)

	if _, err := io.Copy(fileBuffer, file); err != nil {
		return bytes.Buffer{}, err
	}

	return *fileBuffer, nil
}

func isImage(header multipart.FileHeader) error {
	switch header.Header.Get("Content-Type") {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return nil
	default:
		return fmt.Errorf("недопустимый тип файла: %s", header.Header.Get("Content-Type"))
	}
}
