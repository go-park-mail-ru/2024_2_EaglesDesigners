package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/messages/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const pageSize = 25

type MessageRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewMessageRepositoryImpl(pool *pgxpool.Pool) MessageRepository {
	return &MessageRepositoryImpl{
		pool: pool,
	}
}

func (r *MessageRepositoryImpl) GetMessages(page int, chatId uuid.UUID) ([]models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	log.Printf("Repository: соединение успешно установлено")

	rows, err := conn.Query(context.Background(),
		`SELECT
	m.id,
	m.author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	u.username
	FROM public.messages AS m
	JOIN public.user AS u ON u.id = m.author_id
	WHERE m.chat_id = $1
	ORDER BY sent_at DESC
	LIMIT $2
	OFFSET $3;`,
		chatId,
		pageSize,
		pageSize*page,
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
		var authorName string
		var message string
		var sentAt time.Time
		var isRedacted bool

		err = rows.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &authorName)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		messages = append(messages, models.Message{
			MessageId:  messageId,
			AuthorID:   authorID,
			AuthorName: authorName,
			Message:    message,
			SentAt:     sentAt,
			IsRedacted: isRedacted,
		})
	}

	log.Printf("Repository: сообщения успешно найдеты. Количество сообшений: %d", len(messages))
	return messages, nil
}

func (r *MessageRepositoryImpl) AddMessage(message models.Message, chatId uuid.UUID) error {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	log.Printf("Repository: соединение успешно установлено")

	// нужно чё-то придумать со стикерами
	row := conn.QueryRow(context.Background(),
		`INSERT INTO public.message (id, chat_id, author_id, message, sent_at, is_redacted)
	VALUES ($1, $2, $3, $4, $5, false) RETURNING id;`,
		uuid.New(),
		chatId,
		message.AuthorID,
		message.Message,
		message.SentAt,
	)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		log.Printf("Repository: не удалось добавить сообщение: %v", err)
		return err
	}

	return nil
}

func (r *MessageRepositoryImpl) GetLastMessage(chatId uuid.UUID) (models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return models.Message{}, err
	}
	log.Println("Repository: соединение успешно установлено")

	// нужно чё-то придумать со стикерами
	row := conn.QueryRow(context.Background(),
		`SELECT
	m.id,
	m.author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	u.username
	FROM public.message AS m
	JOIN public.user AS u ON u.id = m.author_id
	WHERE m.chat_id = $1
	ORDER BY sent_at DESC
	LIMIT 1;`,
		chatId,
	)

	var messageId uuid.UUID
	var authorID uuid.UUID
	var authorName string
	var message string
	var sentAt time.Time
	var isRedacted bool

	err = row.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &authorName)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Message{}, nil
	}
	if err != nil {
		log.Printf("Repository: unable to scan: %v", err)

		return models.Message{}, err
	}

	messageModel := models.Message{
		MessageId:  messageId,
		AuthorID:   authorID,
		AuthorName: authorName,
		Message:    message,
		SentAt:     sentAt,
		IsRedacted: isRedacted,
	}

	return messageModel, nil
}

func (r *MessageRepositoryImpl) GetAllMessagesAfter(chatId uuid.UUID, after time.Time, lastMessageId uuid.UUID) ([]models.Message, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	log.Printf("Repository: соединение успешно установлено")

	rows, err := conn.Query(context.Background(),
		`SELECT
	m.id,
	m.author_id,
	m.message,
	m.sent_at, 
	m.is_redacted,
	u.username
	FROM public.messages AS m
	JOIN public.user AS u ON u.id = m.author_id
	WHERE m.chat_id = $1 AND m.sent_at >= $2 AND m.id != $3
	ORDER BY sent_at DESC;`,
		chatId,
		after,
		lastMessageId,
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
		var authorName string
		var message string
		var sentAt time.Time
		var isRedacted bool

		err = rows.Scan(&messageId, &authorID, &message, &sentAt, &isRedacted, &authorName)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		messages = append(messages, models.Message{
			MessageId:  messageId,
			AuthorID:   authorID,
			AuthorName: authorName,
			Message:    message,
			SentAt:     sentAt,
			IsRedacted: isRedacted,
		})
	}

	log.Printf("Repository: сообщения успешно найдеты. Количество сообшений: %d", len(messages))
	return messages, nil
}
