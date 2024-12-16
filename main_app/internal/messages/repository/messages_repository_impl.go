package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/messages/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	pageSize                 = 25
	defaultMessageType       = "default"
	informationalMessageType = "informational"
	MessageWithPayloadType   = "with_payload"
	stickerMessageType       = "sticker"
	filePayloadType          = "file"
	photoPayloadType         = "photo"
)

type MessageRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewMessageRepositoryImpl(pool *pgxpool.Pool) MessageRepository {
	return &MessageRepositoryImpl{
		pool: pool,
	}
}

func (r *MessageRepositoryImpl) GetFirstMessages(ctx context.Context, chatId uuid.UUID) ([]models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	defer conn.Release()

	log.Printf("Repository: соединение успешно установлено")

	rows, err := conn.Query(context.Background(),
		`SELECT
	m.id,
	CASE 
        WHEN mt.value = 'informational' THEN '00000000-0000-0000-0000-000000000000' 
        ELSE m.author_id 
    END AS author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	m.branch_id,
	m.chat_id,
	mt.value,
	m.sticker_path
	FROM public.message AS m
	JOIN public.message_type AS mt ON mt.id = m.message_type_id
	WHERE m.chat_id = $1
	ORDER BY sent_at DESC
	LIMIT $2;`,
		chatId,
		pageSize,
	)
	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: сообщения получены")

	messages := []models.Message{}
	for rows.Next() {
		var messageId uuid.UUID
		var authorID uuid.UUID
		var message sql.NullString
		var sentAt time.Time
		var isRedacted bool
		var branchID *uuid.UUID
		var chatID uuid.UUID
		var messageType sql.NullString
		var files []models.Payload
		var photos []models.Payload
		var sticker sql.NullString

		err = rows.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &branchID, &chatID, &messageType, &sticker)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		messages = append(messages, models.Message{
			MessageId:   messageId,
			AuthorID:    authorID,
			Message:     message.String,
			SentAt:      sentAt,
			IsRedacted:  isRedacted,
			BranchID:    branchID,
			ChatId:      chatID,
			MessageType: messageType.String,
			FilesDTO:    files,
			PhotosDTO:   photos,
			Sticker:     sticker.String,
		})
	}

	for i := 0; i < len(messages); i++ {
		if messages[i].MessageType == MessageWithPayloadType {
			log.Printf("поиск вложений сообщения %v", messages[i].MessageId)

			payloadRows, err := conn.Query(context.Background(),
				`select 
					mp.payload_path,
					mp.filename,
					mp.size,
					(SELECT value FROM public.payload_type WHERE id = mp.payload_type) 
				from public.message_payload mp 
				where mp.message_id = $1;`,
				messages[i].MessageId,
			)
			if err != nil {
				log.Printf("Repository: Unable to SELECT payloads: %v\n", err)
				return nil, err
			}

			for payloadRows.Next() {
				var payloadPath string
				var payloadType string
				var size int64
				var filename string

				err = payloadRows.Scan(&payloadPath, &filename, &size, &payloadType)
				if err != nil {
					log.Printf("Repository: unable to scan payloads: %v", err)
					return nil, err
				}
				if payloadType == filePayloadType {
					log.Printf("получен файл %s", payloadPath)
					messages[i].FilesDTO = append(messages[i].FilesDTO, models.Payload{
						URL:      payloadPath,
						Filename: filename,
						Size:     size,
					})
				} else if payloadType == photoPayloadType {
					log.Printf("получено фото %s", payloadPath)
					messages[i].PhotosDTO = append(messages[i].PhotosDTO, models.Payload{
						URL:      payloadPath,
						Filename: filename,
						Size:     size,
					})
				}
			}
		}
	}
	log.Println("Repository: вложения получены")

	log.Printf("Repository: сообщения успешно найдеты. Количество сообшений: %d", len(messages))
	return messages, nil
}

func (r *MessageRepositoryImpl) AddMessage(message models.Message, chatId uuid.UUID) error {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	defer conn.Release()
	log.Printf("Repository: соединение успешно установлено")

	messageType := defaultMessageType

	if message.Sticker != "" {
		messageType = stickerMessageType
	} else if len(message.Files) > 0 || len(message.Photos) > 0 {
		messageType = MessageWithPayloadType
	}

	var row pgx.Row
	if message.Sticker != "" {
		row = conn.QueryRow(context.Background(),
			`INSERT INTO public.message (id, chat_id, author_id, sent_at, is_redacted, message_type_id, sticker_path)
	VALUES ($1, $2, $3, $4, false, (SELECT id FROM message_type WHERE value = $5), $6) RETURNING id;`,
			message.MessageId,
			chatId,
			message.AuthorID,
			message.SentAt,
			messageType,
			message.Sticker,
		)
	} else {
		row = conn.QueryRow(context.Background(),
			`INSERT INTO public.message (id, chat_id, author_id, message, sent_at, is_redacted, message_type_id)
	VALUES ($1, $2, $3, $4, $5, false, (SELECT id FROM message_type WHERE value = $6)) RETURNING id;`,
			message.MessageId,
			chatId,
			message.AuthorID,
			message.Message,
			message.SentAt,
			messageType,
		)
	}

	var message_id uuid.UUID
	if err := row.Scan(&message_id); err != nil {
		log.Printf("Repository: не удалось добавить сообщение: %v", err)
		return err
	}

	if messageType == MessageWithPayloadType {

		for _, file := range message.FilesDTO {
			id := uuid.New()
			row := conn.QueryRow(context.Background(),
				`INSERT INTO public.message_payload (id, message_id, payload_path, filename, size)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
				id,
				message_id,
				file.URL,
				file.Filename,
				file.Size,
			)

			if err := row.Scan(&id); err != nil {
				log.Printf("Repository: не удалось добавить сообщение: %v", err)
				return err
			}
		}

		for _, photo := range message.PhotosDTO {
			id := uuid.New()
			row := conn.QueryRow(context.Background(),
				`INSERT INTO public.message_payload (id, message_id, payload_path, filename, size, payload_type)
		VALUES ($1, $2, $3, $4, $5, (select id from public.payload_type pt where pt.value = $6)) RETURNING id;`,
				id,
				message_id,
				photo.URL,
				photo.Filename,
				photo.Size,
				photoPayloadType,
			)

			if err := row.Scan(&id); err != nil {
				log.Printf("Repository: не удалось добавить сообщение: %v", err)
				return err
			}
		}
	}

	return nil
}

func (r *MessageRepositoryImpl) AddInformationalMessage(message models.Message, chatId uuid.UUID) error {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	defer conn.Release()
	log.Printf("Repository: соединение успешно установлено")

	// нужно чё-то придумать со стикерами
	row := conn.QueryRow(context.Background(),
		`INSERT INTO public.message (id, chat_id, author_id, message, sent_at, is_redacted, message_type_id)
	VALUES ($1, $2, $3, $4, $5, false, (SELECT id FROM message_type WHERE value = 'informational')) RETURNING id;`,
		message.MessageId,
		chatId,
		message.AuthorID,
		message.Message,
		message.SentAt,
	)

	var message_id uuid.UUID
	if err := row.Scan(&message_id); err != nil {
		log.Printf("Repository: не удалось добавить сообщение: %v", err)
		return err
	}

	for _, fileURL := range message.FilesDTO {
		id := uuid.New()
		row := conn.QueryRow(context.Background(),
			`INSERT INTO public.message_payload (id, message_id, payload_path)
		VALUES ($1, $2, $3) RETURNING id;`,
			id,
			message_id,
			fileURL,
		)

		if err := row.Scan(&id); err != nil {
			log.Printf("Repository: не удалось добавить сообщение: %v", err)
			return err
		}
	}

	return nil
}

func (r *MessageRepositoryImpl) DeleteMessage(ctx context.Context, messageId uuid.UUID) error {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`DELETE FROM message WHERE id = $1 RETURNING id`,
		messageId,
	)

	var msgId uuid.UUID
	err = row.Scan(&msgId)

	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepositoryImpl) UpdateMessage(ctx context.Context, messageId uuid.UUID, newText string) error {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`UPDATE message AS m SET
		message = $1,
		is_redacted = true
		WHERE m.id = $2
		RETURNING m.id;`,
		newText,
		messageId,
	)

	var msgId uuid.UUID
	err = row.Scan(&msgId)

	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepositoryImpl) GetMessageById(ctx context.Context, messageId uuid.UUID) (models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return models.Message{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT
		m.author_id,
		m.message,
		m.sent_at, 
		m.is_redacted,
		m.chat_id,
		mt.value
		FROM public.message AS m
		JOIN public.message_type AS mt ON mt.id = m.message_type_id
		WHERE m.id = $1
		ORDER BY sent_at DESC
		LIMIT 1;`,
		messageId,
	)

	var authorID uuid.UUID
	var message string
	var sentAt time.Time
	var isRedacted bool
	var chatId uuid.UUID
	var messageType string

	err = row.Scan(&authorID, &message, &sentAt, &isRedacted, &chatId, &messageType)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Message{}, nil
	}
	if err != nil {
		log.Printf("Repository: unable to scan: %v", err)

		return models.Message{}, err
	}

	messageModel := models.Message{
		MessageId:   messageId,
		AuthorID:    authorID,
		Message:     message,
		SentAt:      sentAt,
		IsRedacted:  isRedacted,
		ChatId:      chatId,
		MessageType: messageType,
	}

	return messageModel, nil
}

func (r *MessageRepositoryImpl) SearchMessagesWithQuery(ctx context.Context, chatId uuid.UUID, searchQuery string) ([]models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT
	m.id,
	m.author_id,
	m.message,
	m.sent_at, 
	m.is_redacted
	FROM public.message AS m
	WHERE m.chat_id = $1 AND lower(m.message) LIKE lower($2)
	ORDER BY sent_at DESC;`,
		chatId,
		"%"+searchQuery+"%",
	)

	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: сообщения получены")

	messages := []models.Message{}
	for rows.Next() {
		var messageId uuid.UUID
		var authorID uuid.UUID
		var message string
		var sentAt time.Time
		var isRedacted bool

		err = rows.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		messages = append(messages, models.Message{
			MessageId:  messageId,
			AuthorID:   authorID,
			Message:    message,
			SentAt:     sentAt,
			IsRedacted: isRedacted,
		})
	}
	log.Printf("Сообщения успешно найдеты. Количество сообшений: %d", len(messages))
	return messages, nil
}

func (r *MessageRepositoryImpl) GetLastMessage(chatId uuid.UUID) (models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return models.Message{}, err
	}
	defer conn.Release()

	// нужно чё-то придумать со стикерами
	row := conn.QueryRow(context.Background(),
		`SELECT
	m.id,
	CASE 
        WHEN mt.value = 'informational' THEN '00000000-0000-0000-0000-000000000000' 
        ELSE m.author_id 
    END AS author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	mt.value
	FROM public.message AS m
	JOIN public.message_type AS mt ON mt.id = m.message_type_id
	WHERE m.chat_id = $1
	ORDER BY sent_at DESC
	LIMIT 1;`,
		chatId,
	)

	var messageId uuid.UUID
	var authorID uuid.UUID
	var message sql.NullString
	var sentAt time.Time
	var isRedacted bool
	var messageType string

	err = row.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &messageType)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Message{}, nil
	}
	if err != nil {
		log.Printf("Repository: unable to scan: %v", err)

		return models.Message{}, err
	}

	messageModel := models.Message{
		MessageId:   messageId,
		AuthorID:    authorID,
		Message:     message.String,
		SentAt:      sentAt,
		IsRedacted:  isRedacted,
		MessageType: messageType,
	}

	if messageModel.MessageType == MessageWithPayloadType {
		log.Printf("поиск вложений сообщения %v", messageModel.MessageId)

		payloadRows, err := conn.Query(context.Background(),
			`select 
				mp.payload_path,
				mp.filename,
				mp.size,
				(SELECT value FROM public.payload_type WHERE id = mp.payload_type) 
			from public.message_payload mp 
			where mp.message_id = $1;`,
			messageModel.MessageId,
		)
		if err != nil {
			log.Printf("Repository: Unable to SELECT payloads: %v\n", err)
			return models.Message{}, err
		}

		for payloadRows.Next() {
			var payloadPath string
			var payloadType string
			var size int64
			var filename string

			err = payloadRows.Scan(&payloadPath, &filename, &size, &payloadType)
			if err != nil {
				log.Printf("Repository: unable to scan payloads: %v", err)
				return models.Message{}, err
			}
			if payloadType == filePayloadType {
				log.Printf("получен файл %s", payloadPath)
				messageModel.FilesDTO = append(messageModel.FilesDTO, models.Payload{
					URL:      payloadPath,
					Filename: filename,
					Size:     size,
				})
			} else if payloadType == photoPayloadType {
				log.Printf("получено фото %s", payloadPath)
				messageModel.PhotosDTO = append(messageModel.PhotosDTO, models.Payload{
					URL:      payloadPath,
					Filename: filename,
					Size:     size,
				})
			}
		}
	}
	log.Println("Repository: вложения получены")

	return messageModel, nil
}

func (r *MessageRepositoryImpl) GetMessagesAfter(ctx context.Context, chatId uuid.UUID, lastMessageId uuid.UUID) ([]models.Message, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	defer conn.Release()
	log.Printf("Repository: соединение успешно установлено")

	rows, err := conn.Query(ctx,
		`SELECT
	m.id,
	CASE 
        WHEN mt.value = 'informational' THEN '00000000-0000-0000-0000-000000000000' 
        ELSE m.author_id 
    END AS author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	mt.value
	FROM public.message AS m
	JOIN public.message_type AS mt ON mt.id = m.message_type_id
	WHERE m.chat_id = $1 AND m.sent_at <= (SELECT sent_at FROM message WHERE id = $2) AND m.id != $2
	ORDER BY sent_at DESC
	LIMIT $3;`,
		chatId,
		lastMessageId,
		pageSize,
	)

	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: сообщения получены")

	messages := []models.Message{}
	for rows.Next() {
		var messageId uuid.UUID
		var authorID uuid.UUID
		var message string
		var sentAt time.Time
		var isRedacted bool
		var messageType string

		err = rows.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &messageType)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		messages = append(messages, models.Message{
			MessageId:   messageId,
			AuthorID:    authorID,
			Message:     message,
			SentAt:      sentAt,
			IsRedacted:  isRedacted,
			MessageType: messageType,
		})
	}

	log.Printf("Repository: сообщения успешно найдеты. Количество сообшений: %d", len(messages))
	return messages, nil
}

func (r *MessageRepositoryImpl) GetPayload(ctx context.Context, chatId uuid.UUID) (files []models.Payload, photos []models.Payload, err error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, nil, err
	}
	defer conn.Release()

	payloadRows, err := conn.Query(context.Background(),
		`SELECT 
			mp.payload_path,
			mp.filename,
			mp.size,
			(SELECT value FROM public.payload_type WHERE id = mp.payload_type) 
		FROM message m
		JOIN message_payload mp ON mp.message_id = m.id
		WHERE m.chat_id = $1;`,
		chatId,
	)
	if err != nil {
		log.Printf("Repository: Unable to SELECT payloads: %v\n", err)
		return nil, nil, err
	}

	for payloadRows.Next() {
		var payloadPath string
		var payloadType string
		var size int64
		var filename string

		err = payloadRows.Scan(&payloadPath, &filename, &size, &payloadType)
		if err != nil {
			log.Printf("Repository: unable to scan payloads: %v", err)
			return nil, nil, err
		}
		if payloadType == filePayloadType {
			log.Printf("получен файл %s", payloadPath)
			files = append(files, models.Payload{
				URL:      payloadPath,
				Filename: filename,
				Size:     size,
			})
		} else if payloadType == photoPayloadType {
			log.Printf("получено фото %s", payloadPath)
			photos = append(photos, models.Payload{
				URL:      payloadPath,
				Filename: filename,
				Size:     size,
			})
		}
	}

	return files, photos, nil
}
