package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
)

func (r *ChatRepositoryImpl) GetSendNotificationsForUser(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (bool, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return false, err
	}
	defer conn.Release()

	var sendNotifications bool

	err = conn.QueryRow(ctx,
		`SELECT 
		send_notifications
		FROM chat_user
		WHERE chat_id = $1 AND user_id = $2`,
		chatId,
		userId,
	).Scan(&sendNotifications)
	if err != nil {
		return false, err
	}
	return sendNotifications, nil
}

func (r *ChatRepositoryImpl) SetChatNotofications(ctx context.Context, chatUUID uuid.UUID, userId uuid.UUID, value bool) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	var sendNotifications bool

	err = conn.QueryRow(ctx,
		`UPDATE chat_user SET
		send_notifications = $1
		WHERE chat_id = $2 AND user_id = $3 RETURNING send_notifications`,
		value,
		chatUUID,
		userId,
	).Scan(&sendNotifications)
	if err != nil {
		return err
	}

	if sendNotifications != value {
		return fmt.Errorf("не удалось обновить значение")
	}

	return nil
}
