package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	userModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/jwt/usecase"
)

type ChatRepositoryImpl struct {
	pool       *pgxpool.Pool
	chat_types map[string]int
}

// readChatTypes подгружает id чатов из бд
func readChatTypes(p *pgxpool.Pool) (map[string]int, error) {
	var chat_types map[string]int

	conn, err := p.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}

	rows, err := conn.Query(context.Background(), "SELECT id, value FROM chat_type;")
	defer conn.Release()

	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		log.Printf("ID: %d, Value: %s\n", id, value)

		chat_types[value] = id
	}
	return chat_types, nil
}

func NewChatRepository(pool *pgxpool.Pool) (ChatRepository, error) {
	chats, err := readChatTypes(pool)
	if err != nil {
		return nil, err
	}

	return &ChatRepositoryImpl{
			pool:       pool,
			chat_types: chats,
		},
		nil
}

func (r *ChatRepositoryImpl) CreateNewChat(chat chatModel.Chat) error {
	chatDAO := chatModel.ChatDAO{
		ChatId:      chat.ChatId,
		ChatName:    chat.ChatName,
		ChatTypeId:  r.chat_types[chat.ChatName],
		AvatarURL:   chat.AvatarURL,
		ChatURLName: chat.ChatURLName,
	}

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO chat (id, chat_name, chat_type_id, avatar_path, chat_link_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`,
		chatDAO.ChatId,
		chatDAO.ChatName,
		chatDAO.ChatTypeId,
		chatDAO.AvatarURL,
		chatDAO.ChatURLName)

	var id uuid.UUID
	err = row.Scan(&id)

	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		return err
	}
	log.Printf("Chat added %s %v", chat.ChatName, chat.ChatId)

	return nil
}

func (r *ChatRepositoryImpl) GetUserChats(user *userModel.User, pageSize int) []chatModel.Chat {
	return nil
}

func (r *ChatRepositoryImpl) IsUserInChat(userId int, chatId int) (bool, error) {
	// идем в бд по двум полям: если есть то тру
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return false, err
	}
	defer conn.Release()

	var id uuid.UUID
	err = conn.QueryRow(context.Background(),
		`SELECT id
		FROM chat_user AS cu
		WHERE cu.user_id = $1 AND cu.chat_id = $2;`,
		userId,
		chatId,
	).Scan(&id)

	if err != nil {
		return false, nil
	}

	return true, nil
}
