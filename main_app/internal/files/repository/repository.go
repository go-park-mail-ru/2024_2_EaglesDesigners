package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/files/models"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	bucket *gridfs.Bucket
	pool   *pgxpool.Pool
}

func New(bucket *gridfs.Bucket, pool *pgxpool.Pool) *Repository {
	return &Repository{
		bucket: bucket,
		pool:   pool,
	}
}

func (r *Repository) GetFile(ctx context.Context, filename string) (*bytes.Buffer, *models.FileMetaData, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	log.Println("поиск по filename:", filename)

	cursor, err := r.bucket.Find(bson.M{"filename": filename})
	if err != nil {
		log.WithError(err).Errorln("не удалось найти файл по имени")
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Err() != nil {
		log.WithError(cursor.Err()).Errorln("Cursor error")
		return nil, nil, cursor.Err()
	}

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

	fileID := fileInfo["_id"].(primitive.ObjectID)

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

	return filename.Hex(), nil
}

func (r *Repository) RewriteFile(ctx context.Context, filename primitive.ObjectID, fileBuffer *bytes.Buffer, metadata primitive.D) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	fileID, err := r.getIDByFilename(ctx, filename)
	if err != nil {
		log.WithError(err).Println("не удалось найти старый файл по fileID")
	} else {
		if err := r.bucket.Delete(fileID); err != nil {
			log.WithError(err).Warnln("старый файл не обнаружен")
		}
	}

	uploadOpts := options.GridFSUpload().SetMetadata(metadata)
	uploadStream, err := r.bucket.OpenUploadStream(filename.Hex(), uploadOpts)
	if err != nil {
		log.WithError(err).Errorln("не удалось открыть поток для загрузки файла")
		return err
	}
	defer func() {
		if err := uploadStream.Close(); err != nil {
			log.WithError(err).Errorln("не удалось закрыть поток")
		}
	}()

	if _, err = uploadStream.Write(fileBuffer.Bytes()); err != nil {
		log.WithError(err).Errorln("не удалось загрузить файл в буфер")
		return err
	}

	return nil
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

func (r *Repository) GetStickerPack(ctx context.Context, packID string) (models.GetStickerPackResponse, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return models.GetStickerPackResponse{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT 
			s.sticker_path,
			sp.photo
		FROM sticker_sticker_pack ssp
		JOIN sticker_pack sp ON ssp.pack = sp.id
		JOIN sticker s ON ssp.sticker = s.id
		WHERE sp.id = $1;`,
		packID,
	)
	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return models.GetStickerPackResponse{}, err
	}
	var pack models.GetStickerPackResponse
	for rows.Next() {
		var url string
		var photo string
		err = rows.Scan(&url, &photo)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return models.GetStickerPackResponse{}, err
		}
		pack.URLs = append(pack.URLs, url)
		pack.Photo = photo
	}

	return pack, nil
}

func (r *Repository) GetStickerPacks(ctx context.Context) (models.StickerPacks, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return models.StickerPacks{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT 
			id,
			photo
		FROM sticker_pack;`,
	)
	if err != nil {
		log.Printf("Repository: Unable to SELECT: %v\n", err)
		return models.StickerPacks{}, err
	}
	var dto models.StickerPacks
	for rows.Next() {
		var id uuid.UUID
		var photo string
		err = rows.Scan(&id, &photo)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return models.StickerPacks{}, err
		}

		dto.Packs = append(dto.Packs, models.StickerPack{
			PackID: id,
			Photo:  photo,
		})
	}

	return dto, nil
}

func (r *Repository) getIDByFilename(ctx context.Context, filename primitive.ObjectID) (primitive.ObjectID, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	cursor, err := r.bucket.Find(bson.M{"filename": filename})
	if err != nil {
		log.WithError(err).Errorln("не удалось найти файл по имени")
		return primitive.ObjectID{}, err
	}
	defer cursor.Close(ctx)

	var fileInfo bson.M
	if cursor.Next(ctx) {
		if err := cursor.Decode(&fileInfo); err != nil {
			log.WithError(err).Errorln("не удалось декодировать метаданные файла")
			return primitive.ObjectID{}, err
		}
	} else {
		log.Errorln("файл не найден")
		return primitive.ObjectID{}, errors.New("файл не найден")
	}

	fileID := fileInfo["_id"].(primitive.ObjectID)

	return fileID, nil
}

func extractID(filePath string) (string, error) {
	// Получаем только имя файла
	fileName := filepath.Base(filePath)

	// Удаляем расширение файла
	id := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return id, nil
}

func (r *Repository) CreateSticker(ctx context.Context, fileBuffer *bytes.Buffer, metadata primitive.D, name string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	fileID, err := extractID(name)
	if err != nil {
		log.WithError(err).Errorln("id err")
		return err
	}

	filename, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		log.WithError(err).Errorln("name err")
		return err
	}

	uploadOpts := options.GridFSUpload().SetMetadata(metadata)

	uploadStream, err := r.bucket.OpenUploadStream(filename.Hex(), uploadOpts)
	if err != nil {
		log.WithError(err).Errorln("не удалось открыть поток для загрузки файла")
		return err
	}
	defer func() {
		if err = uploadStream.Close(); err != nil {
			log.WithError(err).Errorln("во прикол а как так-то")
		}
	}()

	err = uploadStream.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.WithError(err).Errorln("не удалось выставить таймаут")
		return err
	}

	if _, err = uploadStream.Write(fileBuffer.Bytes()); err != nil {
		log.WithError(err).Errorln("не удалось загрузить файл в буфер")
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
