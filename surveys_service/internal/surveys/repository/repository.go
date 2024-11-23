package repository

import (
	"context"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/surveys/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServeyRepository struct {
	pool *pgxpool.Pool
}

func NewServeyRepository(pool *pgxpool.Pool) ServeyRepository {
	return ServeyRepository{
		pool: pool,
	}
}

func (r *ServeyRepository) GetQuestionsByServeyName(ctx context.Context, serveyName string) ([]models.Question, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	defer conn.Release()

	log.Printf("Repository: соединение успешно установлено")

	rows, err := conn.Query(context.Background(),
		`SELECT
	q.id,
	qt.value AS question_type,
	q.question_text
	FROM public.question AS q
	JOIN servey AS s ON s.id = q.servey_id
	JOIN question_type AS qt ON qt.id = q.type_id
	WHERE s.name = $1;`,
		serveyName,
	)
	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: сообщения получены")

	questions := []models.Question{}
	for rows.Next() {
		var questionId uuid.UUID
		var questionText string
		var questionType string

		err = rows.Scan(&questionId, &questionText, &questionType)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		questions = append(questions, models.Question{
			QuestionId:   questionId,
			QuestionText: questionText,
			QuestionType: questionType,
		})
	}

	log.Printf("Repository: сообщения успешно найдеты. Количество сообшений: %d", len(questions))
	return questions, nil
}

// AddAnswer добавляет новый ответ.
func (r *ServeyRepository) AddAnswer(ctx context.Context, answer models.Answer) error {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return err
	}
	defer conn.Release()

	log.Printf("Repository: соединение успешно установлено")

	// нужно чё-то придумать со стикерами
	row := conn.QueryRow(context.Background(),
		`INSERT INTO public.answer (id, question_id, user_id, text_answer, numeric_answer)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		answer.AnswerId,
		answer.QuestionId,
		answer.UserId,
		answer.TextAnswer,
		answer.NumericAnswer,
	)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		log.Printf("Repository: не удалось добавить ответ: %v", err)
		return err
	}

	return nil
}

// GetAllTextAnswers забирает из бд статистику для текстовых вопросов
func (r *ServeyRepository) GetAllTextAnswers(ctx context.Context, questionId uuid.UUID) ([]string, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT 
		a.text_answer
		FROM answer AS a
		WHERE a.question_id = $1;`,
		questionId,
	)

	if err != nil {
		log.Printf("Repository: Unable to SELECT chats: %v\n", err)
		return nil, err
	}
	log.Println("Repository: сообщения получены")

	answers := []string{}
	for rows.Next() {
		var answer string

		err = rows.Scan(&answer)
		if err != nil {
			log.Printf("Repository: unable to scan: %v", err)
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *ServeyRepository) GetAverageNumericAnswer(ctx context.Context, questionId uuid.UUID) (int, error) {
	log := logger.LoggerWithCtx(ctx, logger.Log)

	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Repository: не удалось установить соединение: %v", err)
		return 0, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT AVG(numeric_answer) FROM answer WHERE question_id = $1;`,
		questionId,
	)

	var avg int
	if err := row.Scan(&avg); err != nil {
		log.Printf("Repository: не удалось добавить сообщение: %v", err)
		return 0, err
	}

	return avg, nil
}
