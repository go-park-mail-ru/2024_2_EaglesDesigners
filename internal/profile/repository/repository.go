package repository

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/profile/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (models.ProfileDataDAO, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Не удалось соединиться с базой данных: %v\n", err)
		return models.ProfileDataDAO{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`SELECT 
			birthdate,
			bio,
			avatar_path
		FROM public."user"
		WHERE username = $1;`,
		username,
	)

	var profileData models.ProfileDataDAO

	err = row.Scan(
		&profileData.Birthdate,
		&profileData.Bio,
		&profileData.AvatarURL,
	)
	if err != nil {
		log.Printf("Не удалось получить данные профиля: %v\n", err)
		return models.ProfileDataDAO{}, err
	}

	return profileData, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, profile models.Profile) (avatarURL string, err error) {
	return "", nil
}
