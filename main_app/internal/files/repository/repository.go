package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	bucket *gridfs.Bucket
}

func New(bucket *gridfs.Bucket) *Repository {
	return &Repository{
		bucket: bucket,
	}
}

func (r *Repository) GetFile(ctx context.Context, fileID primitive.ObjectID) (*bytes.Buffer, *models.FileMetaData, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("repo fileID: ", fileID)

	downloadStream, err := r.bucket.OpenDownloadStream(fileID)
	if err != nil {
		log.WithError(err).Errorln("не удалось открыть поток на скачивание файла")
		return nil, nil, err
	}
	defer func() {
		if err := downloadStream.Close(); err != nil {
			log.Panic(err)
		}
	}()

	if err := downloadStream.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		log.WithError(err).Errorln("не удалось выставить таймаут")
		return nil, nil, err
	}

	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, downloadStream); err != nil {
		log.WithError(err).Errorln("не удалось выгрузить файл в буфер")
		return nil, nil, err
	}

	cursor, err := r.bucket.Find(bson.M{"_id": fileID})
	if err != nil {
		log.WithError(err).Errorln("не удалось найти файл по ID")
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var fileInfo bson.M
	if cursor.Next(ctx) {
		if err := cursor.Decode(&fileInfo); err != nil {
			log.WithError(err).Errorln("не удалось декодировать метаданные файла")
			return nil, nil, err
		}
	} else {
		log.Errorln("файл не найден")
		return nil, nil, errors.New("файл не найден")
	}

	metadata := fileInfo["metadata"].(primitive.M)

	log.Print("meta: ", metadata)

	fileMeta := &models.FileMetaData{
		Filename:    metadata["filename"].(string),
		ContentType: metadata["contentType"].(string),
		FileSize:    metadata["size"].(int64),
	}

	log.Print("my meta: ", fileMeta)

	return fileBuffer, fileMeta, nil
}

func (r *Repository) SaveFile(ctx context.Context, fileBuffer *bytes.Buffer, metadata primitive.D) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	
	filename := primitive.NewObjectID()

	uploadOpts := options.GridFSUpload().SetMetadata(metadata)

	uploadStream, err := r.bucket.OpenUploadStream(filename.Hex(), uploadOpts)
	if err != nil {
		log.WithError(err).Errorln("не удалось открыть поток для загрузки файла")
		return "", err
	}
	defer func() {
		if err = uploadStream.Close(); err != nil {
			log.WithError(err).Errorln("во прикол а как так-то")
		}
	}()

	err = uploadStream.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.WithError(err).Errorln("не удалось выставить таймаут")
		return "", err
	}

	if _, err = uploadStream.Write(fileBuffer.Bytes()); err != nil {
		log.WithError(err).Errorln("не удалось загрузить файл в буфер")
		return "", err
	}

	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

func (r *Repository) DeleteFile(ctx context.Context, fileID primitive.ObjectID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	err := r.bucket.Delete(fileID)
	if err != nil {
		log.WithError(err).Errorln("не удалось удалить файл")
		return err
	}

	return nil
}

// func (r *Repository) RewritePhoto(file multipart.File, photoURL string) error {
// 	dst, err := os.Create(photoURL)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, file)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *Repository) RemovePhoto(photoURL string) error {
// 	err := os.Remove(photoURL)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func IsImageFile(file multipart.File) bool {
// 	img, err := imaging.Decode(file)
// 	if err != nil {
// 		return false
// 	}

// 	// сброс указателя на начало файла
// 	if _, err := file.Seek(0, io.SeekStart); err != nil {
// 		return false
// 	}

// 	return img != nil
// }
