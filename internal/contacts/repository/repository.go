package repository

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/internal/contacts/models"
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

func (r *Repository) GetContacts(ctx context.Context, username string) (contacts []models.ContactDAO, err error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Не удалось соединиться с базой данных: %v\n", err)
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
		log.Printf("Repository: Не удалось получить контакты: %v\n", err)
		return contacts, err
	}
	defer rows.Close()

	for rows.Next() {
		var contact models.ContactDAO

		if err = rows.Scan(&contact.ID, &contact.Username, &contact.Name, &contact.AvatarURL); err != nil {
			log.Printf("Repository: Не удалось получить контакты: %v\n", err)
			return contacts, err
		}
		contacts = append(contacts, contact)
	}

	log.Println("Repository: данные получены")

	return contacts, nil
}

func (r *Repository) AddContact(ctx context.Context, contactData models.ContactDataDAO) (models.ContactDAO, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Не удалось соединиться с базой данных: %v\n", err)
		return models.ContactDAO{}, err
	}
	defer conn.Release()

	tx, err := conn.Conn().Begin(ctx)
	if err != nil {
		log.Printf("Repository: Не удалось создать транзацию: %v\n", err)
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
		log.Printf("Repository: Не удалось создать контакт: %v\n", err)
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
		log.Printf("Repository: Не удалось получить данные контакта: %v\n", err)
		return models.ContactDAO{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Repository: Не удалось подтвердить транзакцию: %v\n", err)
		return models.ContactDAO{}, err
	}

	return contact, nil
}
