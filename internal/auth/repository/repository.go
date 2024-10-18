package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db    *pgxpool.Pool
	close func()
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
		close: func() {
			db.Close()
		},
	}
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query := `SELECT 
			  	  u.id,
			  	  u.username,
				  u.password,
				  u.version,
			  	  p.name 
			  FROM public."user" u 
			  JOIN public.profile p ON p.user_id = u.id
			  WHERE username = $1;`

	var user User

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
	)

	if err != nil {
		log.Printf("Пользователь не найден в базе данных: %v\n", err)
		return user, errors.New("user does not exist")
	}

	log.Printf("Пользователь c id %s найден", user.ID.String())
	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, username, name, password string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("Unable to begin transaction: %v\n", err)
		return err
	}

	// rollback если произойдет ошибка
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	firstQuery := `INSERT INTO public.user
				   (
				   	   id,
				   	   username,
				   	   version,
				   	   password
				   ) VALUES ($1, $2, $3, $4) RETURNING id;`

	uuidNew := uuid.New()
	version := 0
	row := tx.QueryRow(
		ctx,
		firstQuery,
		uuidNew,
		username,
		version,
		password,
	)

	var user_id uuid.UUID
	err = row.Scan(&user_id)
	if err != nil {
		log.Printf("Unable to INSERT in TABLE user: %v\n", err)
		return err
	}

	secondQuery := `INSERT INTO public.profile
					(
						id,
						name,
						user_id
					) VALUES ($1, $2, $3);`

	uuidNew = uuid.New()

	_, err = tx.Exec(
		ctx,
		secondQuery,
		uuidNew,
		name,
		user_id,
	)

	if err != nil {
		log.Printf("Unable to INSERT in TABLE profile: %v\n", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Unable to commit transaction: %v\n", err)
		return err
	}

	log.Println("created user:", uuidNew.String(), username)

	return nil
}
