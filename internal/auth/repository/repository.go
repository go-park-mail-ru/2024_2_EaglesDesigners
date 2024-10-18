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
	user, exists := users[username]
	if !exists {
		log.Println("Пользователь не найден в базе данных")
		return user, errors.New("user does not exist")
	}
	log.Printf("Пользователь c id %d найден", user.ID)
	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, username, name, password string) error {
	query := `INSERT INTO public.user
							(
									id
									username,
									version,
									password
							) VALUES ($1, $2, $3, $4) RETURNING id;`

	uuid := uuid.New()
	version := 0
	row := r.db.QueryRow(
		ctx,
		query,
		uuid,
		version,
		username,
		password,
	)

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		return err
	}

	log.Println("created user:", uuid.String(), username)

	return nil
}

var users = map[string]User{
	"user11": {
		Username: "user11",
		Name:     "Бал Матье",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user22": {
		Username: "user22",
		Name:     "Жабка Пепе",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user33": {
		Username: "user33",
		Name:     "Dr Peper",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
	"user44": {
		Username: "user44",
		Name:     "Vincent Vega",
		Password: "e208b28e33d1cb6c69bdddbc5f4298652be5ae2064a8933ce8a97556334715483259a4f4e003c6f5c44a9ceed09b49c792c0a619c5c5a276bbbdcfbd45c6c648",
		Version:  0,
	},
}
