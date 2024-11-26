package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

type Repository struct {
	db    DB
	close func()
}

func NewRepository(db DB) *Repository {
	return &Repository{
		db: db,
		close: func() {
			db.Close()
		},
	}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (models.UserDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	query := `SELECT 
			  	  id,
			  	  username,
				  password,
				  version,
			  	  name,
				  avatar_path
			  FROM public."user"
			  WHERE username = $1;`

	var user models.UserDAO

	row := r.db.QueryRow(
		ctx,
		query,
		username,
	)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Version,
		&user.Name,
		&user.AvatarURL,
	)

	if err != nil {
		log.Errorf("Пользователь не найден в базе данных: %v", err)
		return user, errors.New("пользователь не найден")
	}

	log.Printf("пользователь с id=%s получен", user.ID.String())

	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, username, name, password string) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	query := `INSERT INTO public.user
				   (
				   	   id,
				   	   username,
				   	   version,
				   	   password,
					   name
				   ) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	uuidNew := uuid.New()
	version := 0
	row := r.db.QueryRow(
		ctx,
		query,
		uuidNew,
		username,
		version,
		password,
		name,
	)

	var user_id uuid.UUID
	err := row.Scan(&user_id)
	if err != nil {
		log.Errorf("не удалось вставить в таблицу юзера: %v", err)
		return err
	}

	log.Println("пользователь успешно создан: ", uuidNew.String(), username)

	return nil
}
