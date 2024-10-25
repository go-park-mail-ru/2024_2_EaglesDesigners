package repository

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
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

func (r *Repository) GetContacts(ctx context.Context, username string) (contacts []models.UserDAO, err error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Не удалось соединиться с базой данных: %v\n", err)
		return contacts, err
	}
	defer conn.Release()

	log.Println("repo before query")

	rows, err := conn.Query(
		ctx,
		`SELECT 
			username,
			name,
			avatar_path
		FROM public."user"
		WHERE id IN 
		(
			SELECT contact_id 
			FROM public."contact"
			WHERE user_id = (SELECT id FROM public."user" WHERE username = $1)
		);`,
		username,
	)
	if err != nil {
		log.Printf("Не удалось получить контакты: %v\n", err)
		return contacts, err
	}
	defer rows.Close()

	log.Println("repo after query")

	for rows.Next() {
		var user models.UserDAO

		if err = rows.Scan(&user.Username, &user.Name, &user.AvatarURL); err != nil {
			log.Printf("Не удалось получить контакты: %v\n", err)
			return contacts, err
		}
		contacts = append(contacts, user)
	}

	log.Println("repo cont", contacts)

	log.Println("repo done")

	return contacts, nil
}

// UPDATE public."user"
// SET
// name = $2,
// bio = $3,
// birthdate = $4,
// avatar_path = $5
// WHERE username = $1
// RETURNING avatar_path;
// func (r *Repository) UpdateProfile(ctx context.Context, profile models.Profile) (avatarURL *string, err error) {
// 	conn, err := r.pool.Acquire(ctx)
// 	if err != nil {
// 		log.Printf("Не удалось соединиться с базой данных: %v\n", err)
// 		return nil, err
// 	}
// 	defer conn.Release()

// 	row := conn.QueryRow(
// 		ctx,
// 		`SELECT avatar_path
// 		FROM public."user"
// 		WHERE username = $1;`,
// 		profile.Username,
// 	)

// 	err = row.Scan(&avatarURL)
// 	if err != nil {
// 		return nil, errors.New("не получилось получить avatarURL")
// 	}

// 	if avatarURL == nil {
// 		avatarURL = new(string)
// 		*avatarURL = uuid.New().String()
// 	}

// 	query := `UPDATE public."user" SET `
// 	var rowsWithFields []string

// 	var args []interface{}

// 	args = append(args, profile.Username)

// 	if profile.Name != nil {
// 		rowsWithFields = append(rowsWithFields, fmt.Sprintf("name = $%d", len(args)+1))
// 		args = append(args, profile.Name)
// 	}
// 	if profile.Bio != nil {
// 		rowsWithFields = append(rowsWithFields, fmt.Sprintf("bio = $%d", len(args)+1))
// 		args = append(args, profile.Bio)
// 	}
// 	if profile.AvatarBase64 != nil {
// 		rowsWithFields = append(rowsWithFields, fmt.Sprintf("avatar_path = $%d", len(args)+1))
// 		args = append(args, avatarURL)
// 	}
// 	if profile.Birthdate != nil {
// 		rowsWithFields = append(rowsWithFields, fmt.Sprintf("birthdate = $%d", len(args)+1))
// 		args = append(args, profile.Birthdate)
// 	}

// 	if len(args) == 1 {
// 		return nil, errors.New("нет полей для обновления")
// 	}

// 	query += fmt.Sprintf("%s WHERE username = $1", strings.Join(rowsWithFields, ", ")) + " RETURNING avatar_path;"

// 	log.Println(query)

// 	row = conn.QueryRow(ctx, query, args...)

// 	err = row.Scan(&avatarURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return avatarURL, nil
// }
