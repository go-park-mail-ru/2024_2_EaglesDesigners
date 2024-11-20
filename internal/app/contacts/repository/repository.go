package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrContactAlreadyExist = errors.New("contact already exist")

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetContacts(ctx context.Context, username string) (contacts []models.ContactDAO, err error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Не удалось соединиться с базой данных: %v", err)
		return contacts, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			id,
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
		log.Errorf("Не удалось получить контакты: %v", err)
		return contacts, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactDAO

		if err = rows.Scan(&contact.ID, &contact.Username, &contact.Name, &contact.AvatarURL); err != nil {
			log.Errorf("Не удалось получить контакты: %v", err)
			return contacts, err
		}
		contacts = append(contacts, contact)
	}

	log.Println("данные получены")

	return contacts, nil
}

func (r *Repository) AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Не удалось соединиться с базой данных: %v", err)
		return models.ContactDAO{}, err
	}
	defer conn.Release()

	tx, err := conn.Conn().Begin(ctx)
	if err != nil {
		log.Errorf("Не удалось создать транзацию: %v", err)
		return models.ContactDAO{}, err
	}

	newUUID := uuid.New()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO public.contact 
		(id, 
		user_id, 
		contact_id
		)
		VALUES ($1,$2, (SELECT id FROM public."user" WHERE username = $3));`,
		newUUID.String(),
		contactData.UserID,
		contactData.ContactUsername,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Errorf("Контакт уже существует")
			return models.ContactDAO{}, ErrContactAlreadyExist
		}

		log.Errorf("Не удалось создать контакт: %v", err)
		return models.ContactDAO{}, err
	}

	var contact models.ContactDAO

	contact.Username = contactData.ContactUsername

	err = tx.QueryRow(
		ctx,
		`SELECT 
			id,
			name,
			avatar_path
		FROM public."user"
		WHERE username = $1;`,
		contactData.ContactUsername,
	).Scan(&contact.ID, &contact.Name, &contact.AvatarURL)

	if err != nil {
		log.Errorf("Не удалось получить данные контакта: %v", err)
		return models.ContactDAO{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("Не удалось подтвердить транзакцию: %v", err)
		return models.ContactDAO{}, err
	}

	return contact, nil
}

func (r *Repository) DeleteContact(ctx context.Context, contactData models.ContactDataDAO) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Не удалось соединиться с базой данных: %v", err)
		return err
	}
	defer conn.Release()

	result, err := conn.Exec(
		ctx,
		`DELETE FROM public.contact 
		WHERE user_id = $1 AND contact_id = (SELECT id FROM public."user" WHERE username = $2);`,
		contactData.UserID,
		contactData.ContactUsername,
	)
	if err != nil {
		log.Errorf("Не удалось удалить контакт: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		log.Errorf("ничего не удалено; возможно, контакт не найден")
		return errors.New("контакт не найден")
	} else {
		log.Println("контакт удален")
	}

	return nil
}

func (r *Repository) SearchUserContacts(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Не удалось соединиться с базой данных: %v", err)
		return nil, err
	}
	defer conn.Release()

	var contacts []models.ContactDAO

	rows, err := conn.Query(
		ctx,
		`SELECT
			id,
			username,
			name,
			avatar_path
		FROM public."user"
		WHERE id IN 
		(
			SELECT contact_id 
			FROM public."contact"
			WHERE 
				user_id = $1 AND 
				(POSITION(LOWER($2) IN LOWER(username)) > 0 OR POSITION(LOWER($2) IN LOWER(name)) > 0)
		);`,
		userID,
		keyWord,
	)
	if err != nil {
		log.Errorf("Не удалось получить контакты: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactDAO

		if err = rows.Scan(&contact.ID, &contact.Username, &contact.Name, &contact.AvatarURL); err != nil {
			log.Errorf("Не удалось получить контакты: %v", err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	log.Debugln("данные получены")

	return contacts, nil
}

func (r *Repository) SearchGlobalUsers(ctx context.Context, userID uuid.UUID, keyWord string) ([]models.ContactDAO, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Не удалось соединиться с базой данных: %v", err)
		return nil, err
	}
	defer conn.Release()

	var contacts []models.ContactDAO

	rows, err := conn.Query(
		ctx,
		`SELECT
			u.id,
			u.username,
			u.name,
			u.avatar_path
		FROM (
			SELECT
				id,
				username,
				name,
				avatar_path
			FROM public."user"
			WHERE 
				id <> $1 AND
				(POSITION(LOWER($2) IN LOWER(username)) > 0 OR POSITION(LOWER($2) IN LOWER(name)) > 0)
		) AS u
		WHERE id NOT IN (
			SELECT contact_id 
			FROM public."contact"
			WHERE user_id = $1
		);`,
		userID,
		keyWord,
	)
	if err != nil {
		log.Errorf("Не удалось получить контакты: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactDAO

		if err = rows.Scan(&contact.ID, &contact.Username, &contact.Name, &contact.AvatarURL); err != nil {
			log.Errorf("Не удалось получить контакты: %v", err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	log.Debugln("данные получены")

	return contacts, nil
}
