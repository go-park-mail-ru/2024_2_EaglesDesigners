package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/chats/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
)

type ChatRepositoryImpl struct {
	pool       *pgxpool.Pool
	chat_types map[string]int
}

const pageSize = 25

// readChatTypes подгружает id чатов из бд
func readChatTypes(p *pgxpool.Pool) (map[string]int, error) {
	var chat_types map[string]int = map[string]int{}

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

func (r *ChatRepositoryImpl) CreateNewChat(ctx context.Context, chat chatModel.Chat) error {
	chatDAO := chatModel.ChatDAO{
		ChatId:      chat.ChatId,
		ChatName:    chat.ChatName,
		ChatTypeId:  r.chat_types[chat.ChatType],
		AvatarURL:   chat.AvatarURL,
		ChatURLName: chat.ChatURLName,
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	var row pgx.Row
	if chatDAO.ChatURLName == "" && chatDAO.AvatarURL == "" {
		row = conn.QueryRow(ctx,
			`INSERT INTO chat (id, chat_name, chat_type_id)
		VALUES ($1, $2, $3)
		RETURNING id;`,
			chatDAO.ChatId,
			chatDAO.ChatName,
			chatDAO.ChatTypeId,
		)
	} else if chatDAO.ChatURLName == "" {
		row = conn.QueryRow(ctx,
			`INSERT INTO chat (id, chat_name, chat_type_id, avatar_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`,
			chatDAO.ChatId,
			chatDAO.ChatName,
			chatDAO.ChatTypeId,
			chatDAO.AvatarURL,
		)
	} else if chatDAO.AvatarURL == "" {
		row = conn.QueryRow(ctx,
			`INSERT INTO chat (id, chat_name, chat_type_id, chat_link_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`,
			chatDAO.ChatId,
			chatDAO.ChatName,
			chatDAO.ChatTypeId,
			chatDAO.ChatURLName,
		)
	} else {
		row = conn.QueryRow(ctx,
			`INSERT INTO chat (id, chat_name, chat_type_id, avatar_path, chat_link_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`,
			chatDAO.ChatId,
			chatDAO.ChatName,
			chatDAO.ChatTypeId,
			chatDAO.AvatarURL,
			chatDAO.ChatURLName,
		)
	}
	var id uuid.UUID
	err = row.Scan(&id)

	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		return err
	}
	log.Printf("Chat added %s %v", chat.ChatName, chat.ChatId)

	return nil
}

func (r *ChatRepositoryImpl) GetUserChats(ctx context.Context, userId uuid.UUID, pageNum int) ([]chatModel.Chat, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	log.Println("Repository: Соединение с бд установлено")

	rows, err := conn.Query(ctx,
		`SELECT c.id,
		c.chat_name,
		ch.value,
		c.avatar_path,
		c.chat_link_name
		FROM chat_user AS cu
		JOIN chat AS c ON c.id = cu.chat_id
		JOIN chat_type AS ch ON ch.id = c.chat_type_id
		WHERE cu.user_id = $1
		LIMIT $2
		OFFSET $3;`,
		userId,
		pageSize,
		pageSize*pageNum,
	)
	if err != nil {
		log.Printf("Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: Чаты получены")

	chats := []chatModel.Chat{}
	for rows.Next() {
		var chatId uuid.UUID
		var chatName string
		var chatType string
		var avatarURL sql.NullString
		var chatURLName sql.NullString

		log.Println("Repository: поиск параметров из запроса")
		err = rows.Scan(&chatId, &chatName, &chatType, &avatarURL, &chatURLName)

		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		chats = append(chats, chatModel.Chat{
			ChatId:      chatId,
			ChatName:    chatName,
			ChatType:    chatType,
			AvatarURL:   avatarURL.String,
			ChatURLName: chatURLName.String,
		})
	}

	log.Printf("Repository: чаты успешно найдеты. Количество чатов: %d", len(chats))
	return chats, nil
}

func (r *ChatRepositoryImpl) GetUserRoleInChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	// идем в бд по двум полям: если есть то тру
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return "", err
	}
	defer conn.Release()

	var role string
	err = conn.QueryRow(ctx,
		`SELECT ur.value
		FROM chat_user AS cu
		JOIN user_role AS ur ON ur.id = cu.user_role_id
		WHERE cu.user_id = $1 AND cu.chat_id = $2;`,
		userId,
		chatId,
	).Scan(&role)

	if err != nil {
		return "", nil
	}

	return role, nil
}

func (r *ChatRepositoryImpl) AddUserIntoChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID, userROle string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	var id uuid.UUID
	err = conn.QueryRow(ctx,
		`INSERT INTO chat_user (id, user_role_id, chat_id, user_id)
		VALUES ($1, (SELECT id FROM user_role WHERE value = $2), $3, $4)
		RETURNING id;`,
		uuid.New(),
		userROle,
		chatId,
		userId,
	).Scan(&id)

	if err != nil {
		log.Printf("польтзователь %v не добавлен в чат %v. Ошибка: ", userId, chatId, err)
		return err
	}
	log.Printf("польтзователь %v добавлен в чат %v", userId, chatId)
	return nil
}

func (r *ChatRepositoryImpl) GetCountOfUsersInChat(ctx context.Context, chatId uuid.UUID) (int, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	log.Printf("Chat repository: установка количества участников чата: %v", chatId)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return 0, err
	}
	defer func() { conn.Release() }()

	var count int

	err = conn.QueryRow(ctx,
		`SELECT COUNT(id)
		FROM chat_user AS cu
		WHERE cu.chat_id = $1;`,
		chatId,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, err
}

func (r *ChatRepositoryImpl) GetChatById(ctx context.Context, chatId uuid.UUID) (chatModel.Chat, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return chatModel.Chat{}, err
	}
	defer conn.Release()

	var chatName string
	var chatType string
	var avatarURL sql.NullString
	var chatURLName sql.NullString

	err = conn.QueryRow(ctx,
		`SELECT c.id,
		c.chat_name,
		ch.value,
		c.avatar_path,
		c.chat_link_name
		FROM chat AS c
		JOIN chat_type AS ch ON ch.id = c.chat_type_id
		WHERE c.id = $1`,
		chatId,
	).Scan(&chatId, &chatName, &chatType, &avatarURL, &chatURLName)

	if err != nil {
		return chatModel.Chat{}, nil
	}

	return chatModel.Chat{
		ChatId:      chatId,
		ChatName:    chatName,
		ChatType:    chatType,
		AvatarURL:   avatarURL.String,
		ChatURLName: chatURLName.String,
	}, nil

}

func (r *ChatRepositoryImpl) DeleteChat(ctx context.Context, chatId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	log.Printf("Chat repository -> DeleteChat: начато удаление чата: %v", chatId)

	deleteQuery := `DELETE FROM chat WHERE id = $1;`

	// Выполнение удаления
	_, err = conn.Exec(context.Background(), deleteQuery, chatId)

	if err != nil {
		log.Printf("Chat repository -> DeleteChat: не удалось удалить чат: %v", err)
		return err
	}

	return nil
}

func (r *ChatRepositoryImpl) UpdateChat(ctx context.Context, chatId uuid.UUID, chatName string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	log.Printf("Chat repository -> UpdateChat: начато обновление чата: %v", chatId)

	deleteQuery := `UPDATE chat SET
		chat_name = $1 WHERE id = $2;`

	// Выполнение удаления
	_, err = conn.Exec(ctx, deleteQuery, chatName, chatId)

	if err != nil {
		log.Printf("Chat repository -> UpdateChat: не удалось обновить чат: %v", err)
		return err
	}

	return nil
}

func (r *ChatRepositoryImpl) DeleteUserFromChat(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	log.Printf("Chat repository -> UpdateChat: начато обновление чата: %v", chatId)

	deleteQuery := `DELETE FROM chat_user WHERE chat_id = $1 AND user_id = $2;`

	// Выполнение удаления
	_, err = conn.Exec(ctx, deleteQuery, chatId, userId)

	if err != nil {
		log.Printf("Chat repository -> DeleteUserFromChat: не удалось обновить чат: %v", err)
		return err
	}

	return nil
}

func (r *ChatRepositoryImpl) GetUsersFromChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()

	log.Printf("Chat repository -> GetUsersFromChat: начато получение пользователей из чата: %v", chatId)

	rows, err := conn.Query(ctx,
		`SELECT user_id
		FROM chat_user
		WHERE chat_id = $1;`,
		chatId,
	)
	if err != nil {
		log.Printf("Unable to SELECT ids: %v\n", err)
		return nil, err
	}

	var ids []uuid.UUID
	for rows.Next() {
		var userId uuid.UUID

		log.Println("Repository: поиск параметров из запроса")
		err = rows.Scan(&userId)

		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}
		ids = append(ids, userId)
	}
	return ids, nil
}

func (r *ChatRepositoryImpl) UpdateChatPhoto(ctx context.Context, chatId uuid.UUID, filename string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	log.Printf("Chat repository -> UpdateChatPhoto: начато обновление чата: %v", chatId)

	deleteQuery := `UPDATE chat SET
		avatar_path = $1 WHERE id = $2;`

	// Выполнение удаления
	_, err = conn.Exec(ctx, deleteQuery, filename, chatId)

	if err != nil {
		log.Printf("Chat repository -> UpdateChatPhoto: не удалось обновить чат: %v", err)
		return err
	}

	return nil
}

func (r *ChatRepositoryImpl) GetUserNameAndAvatar(ctx context.Context, userId uuid.UUID) (string, string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return "", "", err
	}
	defer conn.Release()

	var username sql.NullString
	var filename sql.NullString
	err = conn.QueryRow(ctx,
		`SELECT username, avatar_path FROM public.user WHERE id = $1;`,
		userId,
	).Scan(&username, &filename)

	if err != nil {
		log.Printf("Chat repository -> GetUserNameAndAvatar: не удалось получитьб юзера: %v", err)
		return "", "", err
	}

	return username.String, filename.String, nil
}
