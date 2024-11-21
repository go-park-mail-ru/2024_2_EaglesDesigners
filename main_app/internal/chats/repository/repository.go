package repository

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	errGroup "golang.org/x/sync/errgroup"

	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/logger"
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

func (r *ChatRepositoryImpl) GetUserChats(ctx context.Context, userId uuid.UUID) ([]chatModel.Chat, error) {
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
		WHERE cu.user_id = $1;`,
		userId,
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

func (r *ChatRepositoryImpl) GetChatType(ctx context.Context, chatId uuid.UUID) (string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Repository: Unable to acquire a database connection: %v", err)
		return "", err
	}
	defer conn.Release()

	var chatType string

	err = conn.QueryRow(ctx,
		`SELECT ct.value 
		FROM chat ch
		JOIN chat_type ct ON ct.id = ch.chat_type_id 
		WHERE ch.id = $1;`,
		chatId,
	).Scan(&chatType)

	if err != nil {
		log.Errorf("Не удалось найти тип чата: %v", err)
		return "", err
	}

	return chatType, nil
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

func (r *ChatRepositoryImpl) GetUsersFromChat(ctx context.Context, chatId uuid.UUID) ([]chatModel.UserInChatDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	log.Printf("начато получение пользователей из чата: %v", chatId)

	rows, err := conn.Query(ctx,
		`SELECT 
			u.id,
			u.username,
			u.name,
			u.avatar_path,
			ch.user_role_id
		FROM public.chat_user AS ch
		JOIN public."user" u ON ch.user_id = u.id 
		WHERE chat_id = $1;`,
		chatId,
	)
	if err != nil {
		log.Printf("Unable to SELECT ids: %v", err)
		return nil, err
	}

	var users []chatModel.UserInChatDAO

	log.Println("поиск параметров из запроса")

	var mu sync.Mutex
	var g errGroup.Group

	for rows.Next() {
		g.Go(func() error {
			var user chatModel.UserInChatDAO

			err = rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarPath, &user.Role)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			users = append(users, user)
			return nil
		})
		if err := g.Wait(); err != nil {
			log.Printf("unable to scan: %v", err)
			return []chatModel.UserInChatDAO{}, err
		}
	}

	g.Wait()

	return users, nil
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

func (r *ChatRepositoryImpl) GetNameAndAvatar(ctx context.Context, userId uuid.UUID) (string, string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return "", "", err
	}
	defer conn.Release()

	var name sql.NullString
	var filename sql.NullString
	err = conn.QueryRow(ctx,
		`SELECT name, avatar_path FROM public.user WHERE id = $1;`,
		userId,
	).Scan(&name, &filename)

	if err != nil {
		log.Printf("Chat repository -> GetNameAndAvatar: не удалось получитьб юзера: %v", err)
		return "", "", err
	}

	return name.String, filename.String, nil
}

func (r *ChatRepositoryImpl) AddBranch(ctx context.Context, chatId uuid.UUID, messageID uuid.UUID) (chatModel.AddBranch, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return chatModel.AddBranch{}, err
	}
	defer conn.Release()

	tx, err := conn.Conn().Begin(ctx)
	if err != nil {
		log.Printf("Repository: Unable to create transaction: %v\n", err)
		return chatModel.AddBranch{}, err
	}

	var branch chatModel.AddBranch
	branch.ID = uuid.New()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO public.chat 
		(id,
		chat_name,
		chat_type_id
		)
		VALUES ($1, 'branch', (SELECT id FROM public.chat_type WHERE value = 'branch'))`,
		branch.ID,
	)

	if err != nil {
		log.Errorf("Не удалось добавить ветку: %v", err)
		return chatModel.AddBranch{}, err
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE public.message 
		SET branch_id = $2
		WHERE id = $1;`,
		messageID,
		branch.ID,
	)

	if err != nil {
		log.Errorf("Не удалось привязать ветку к сообщению: %v", err)
		return chatModel.AddBranch{}, err
	}

	log.Debugf("вставка юзеров в ветку %s чата %s", branch.ID.String(), chatId)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO public.chat_user 
			(id, 
			user_role_id, 
			chat_id, 
			user_id)
		SELECT 
			gen_random_uuid(),
			(SELECT id FROM public.user_role WHERE value = 'none'), 
			$2, 
			cu.user_id 
		FROM public.chat_user cu
		WHERE cu.chat_id = $1;`,
		chatId,
		branch.ID,
	)

	if err != nil {
		log.Errorf("Не удалось добавить пользователей в ветку: %v", err)
		return chatModel.AddBranch{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("Не удалось подтвердить транзакцию: %v", err)
		return chatModel.AddBranch{}, err
	}

	return branch, nil
}

func (r *ChatRepositoryImpl) SearchUserChats(ctx context.Context, userId uuid.UUID, keyWord string) ([]chatModel.Chat, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	log.Debugln("Соединение с бд установлено")

	rows, err := conn.Query(ctx,
		`SELECT c.id,
			c.chat_name,
			ch.value,
			c.avatar_path,
			c.chat_link_name
		FROM chat_user AS cu
		JOIN chat AS c ON c.id = cu.chat_id
		JOIN chat_type AS ch ON ch.id = c.chat_type_id
		WHERE 
			cu.user_id = $1 AND 
			ch.value <> 'branch' AND 
			(POSITION(LOWER($2) IN LOWER(c.chat_name)) > 0 OR POSITION(LOWER($2) IN LOWER(c.chat_link_name)) > 0);`,
		userId,
		keyWord,
	)
	if err != nil {
		log.Errorf("Unable to SELECT chats: %v", err)
		return nil, err
	}
	log.Debugln("Чаты получены")

	chats := []chatModel.Chat{}
	for rows.Next() {
		var chatId uuid.UUID
		var chatName string
		var chatType string
		var avatarURL sql.NullString
		var chatURLName sql.NullString

		log.Debugln("поиск параметров из запроса")
		err = rows.Scan(&chatId, &chatName, &chatType, &avatarURL, &chatURLName)

		if err != nil {
			log.Errorf("unable to scan: %v", err)
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

	log.Debugf("чаты успешно найдеты. Количество чатов: %d", len(chats))

	return chats, nil
}

func (r *ChatRepositoryImpl) SearchGlobalChats(ctx context.Context, userId uuid.UUID, keyWord string) ([]chatModel.Chat, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	log.Debugln("Соединение с бд установлено")

	rows, err := conn.Query(ctx,
		`SELECT 
			ch.id,
			ch.chat_name,
			ch.value,
			ch.avatar_path,
			ch.chat_link_name
		FROM (
			SELECT 
				c.id,
				c.chat_name,
				ch.value,
				c.avatar_path,
				c.chat_link_name
			FROM public.chat c
			JOIN public.chat_type ch ON ch.id = c.chat_type_id
			WHERE  
				ch.value = 'channel' AND 
				(POSITION(LOWER($2) IN LOWER(c.chat_name)) > 0 OR POSITION(LOWER($2) IN LOWER(c.chat_link_name)) > 0)
		) AS ch
		WHERE ch.id NOT IN (
			SELECT c.id
			FROM public.chat_user cu
			JOIN public.chat c ON cu.chat_id = c.id 
			WHERE cu.user_id = $1
		);`,
		userId,
		keyWord,
	)
	if err != nil {
		log.Errorf("Unable to SELECT chats: %v", err)
		return nil, err
	}
	log.Debugln("Чаты получены")

	chats := []chatModel.Chat{}
	for rows.Next() {
		var chatId uuid.UUID
		var chatName string
		var chatType string
		var avatarURL sql.NullString
		var chatURLName sql.NullString

		log.Debugln("поиск параметров из запроса")
		err = rows.Scan(&chatId, &chatName, &chatType, &avatarURL, &chatURLName)

		if err != nil {
			log.Errorf("unable to scan: %v", err)
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

	log.Debugf("чаты успешно найдеты. Количество чатов: %d", len(chats))

	return chats, nil
}
