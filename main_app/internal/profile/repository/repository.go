package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/profile/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/logger"
	"github.com/google/uuid"
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

func (r *Repository) GetProfileByUsername(ctx context.Context, id uuid.UUID) (models.ProfileDataDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("не удалось соединиться с базой данных: %v\n", err)
		return models.ProfileDataDAO{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(
		ctx,
		`SELECT
			name, 
			birthdate,
			bio,
			avatar_path
		FROM public."user"
		WHERE id = $1;`,
		id,
	)

	var profileData models.ProfileDataDAO

	err = row.Scan(
		&profileData.Name,
		&profileData.Birthdate,
		&profileData.Bio,
		&profileData.AvatarPath,
	)
	if err != nil {
		log.Errorf("не удалось получить данные профиля: %v\n", err)
		return models.ProfileDataDAO{}, err
	}

	log.Println("данные получены")

	return profileData, nil
}

// UPDATE public."user"
// SET
// name = $2,
// bio = $3,
// birthdate = $4,
// avatar_path = $5
// WHERE id = $1
// RETURNING avatar_path;
func (r *Repository) UpdateProfile(ctx context.Context, profile models.Profile) (avatarNewURL *string, avatarOldURL *string, err error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("не удалось соединиться с базой данных: %v\n", err)
		return nil, nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(
		ctx,
		`SELECT avatar_path
		FROM public."user"
		WHERE id = $1;`,
		profile.ID,
	)

	err = row.Scan(&avatarOldURL)
	if err != nil {
		return nil, nil, errors.New("не получилось получить avatarURL")
	}

	query := `UPDATE public."user" SET `
	var rowsWithFields []string

	var args []interface{}

	args = append(args, profile.ID)

	if profile.Name != nil {
		rowsWithFields = append(rowsWithFields, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, profile.Name)
	}
	if profile.Bio != nil {
		rowsWithFields = append(rowsWithFields, fmt.Sprintf("bio = $%d", len(args)+1))
		args = append(args, profile.Bio)
	}
	if profile.Avatar != nil {
		rowsWithFields = append(rowsWithFields, fmt.Sprintf("avatar_path = $%d", len(args)+1))

		avatarNewURL = new(string)
		*avatarNewURL = "/uploads/avatar/" + uuid.New().String() + ".png"

		args = append(args, avatarNewURL)
	}
	if profile.DeleteAvatar {
		rowsWithFields = append(rowsWithFields, fmt.Sprintf("avatar_path = $%d", len(args)+1))

		avatarNewURL = new(string)
		avatarNewURL = nil

		args = append(args, avatarNewURL)
	}
	if profile.Birthdate != nil {
		rowsWithFields = append(rowsWithFields, fmt.Sprintf("birthdate = $%d", len(args)+1))
		args = append(args, profile.Birthdate)
	}

	if len(args) == 1 {
		return nil, nil, errors.New("нет полей для обновления")
	}

	query += fmt.Sprintf("%s WHERE id = $1", strings.Join(rowsWithFields, ", ")) + " RETURNING avatar_path;"

	log.Println(query)

	row = conn.QueryRow(ctx, query, args...)

	err = row.Scan(&avatarNewURL)
	if err != nil {
		return nil, nil, err
	}

	return avatarNewURL, avatarOldURL, nil
}
