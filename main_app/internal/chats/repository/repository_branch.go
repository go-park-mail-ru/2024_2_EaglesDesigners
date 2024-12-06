package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	chatModel "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/chats/models"
	"github.com/google/uuid"
)

func (r *ChatRepositoryImpl) AddBranch(ctx context.Context, branchId uuid.UUID, messageID uuid.UUID) (chatModel.AddBranch, error) {
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
	branch.ID = messageID

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

	log.Debugf("вставка юзеров в ветку %s чата %s", branch.ID.String(), branchId)

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
		branchId,
		branch.ID,
	)

	if err != nil {
		log.Errorf("Не удалось добавить пользователей в ветку: %v", err)
		return chatModel.AddBranch{}, err
	}

	parentId, err := r.GetBranchParent(ctx, branchId)
	if err != nil {
		log.Errorf("Не удалось получить родителя: %v", err)
		return chatModel.AddBranch{}, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO chat_brahch ($1, $2, $3)`,
		uuid.New(),
		parentId,
		branchId,
	)
	if err != nil {
		log.Printf("Не удослоь добавить чату ветку: %v", err)
		return chatModel.AddBranch{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Errorf("Не удалось подтвердить транзакцию: %v", err)
		return chatModel.AddBranch{}, err
	}

	return branch, nil
}

// GetBranchParent находит родительский чат. Рассчет идет из того, что branchId == chatId.
func (r *ChatRepositoryImpl) GetBranchParent(ctx context.Context, branchId uuid.UUID) (uuid.UUID, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		log.Printf("Repository: Unable to acquire a database connection: %v\n", err)
		return uuid.Nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		`SELECT
		m.chat_id
		FROM message AS m WHERE m.id = $1;`,
		branchId,
	)

	var chatId uuid.UUID

	row.Scan(&chatId)

	return chatId, nil
}
