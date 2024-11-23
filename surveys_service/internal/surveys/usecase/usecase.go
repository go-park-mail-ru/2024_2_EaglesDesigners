package usecase

import (
	"context"
	"fmt"

	surveysv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/proto"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/surveys/models"

	"github.com/google/uuid"
)

type ServeyRepository interface {
	GetQuestionsByServeyName(ctx context.Context, serveyName string) ([]models.Question, error)
	GetServey(ctx context.Context, serveyName string) (models.Servey, error)
	AddAnswer(ctx context.Context, answer models.Answer) error
	GetAllTextAnswers(ctx context.Context, questionId uuid.UUID) ([]string, error)
	GetAverageNumericAnswer(ctx context.Context, questionId uuid.UUID) (int, error)
	GetQuestionType(ctx context.Context, questionId uuid.UUID) (string, error)
}

type ServeyUsecase struct {
	repository ServeyRepository
}

func NewUsecase(repository ServeyRepository) ServeyUsecase {
	return ServeyUsecase{
		repository: repository,
	}
}

func (u ServeyUsecase) GetSurvey(ctx context.Context, in *surveysv1.GetSurveyReq) (*surveysv1.GetSurveyResp, error) {
	serveyName := in.Name
	servey, err := u.repository.GetServey(ctx, serveyName)
	if err != nil {
		return nil, err
	}

	questions, err := u.repository.GetQuestionsByServeyName(ctx, serveyName)
	if err != nil {
		return nil, err
	}

	questionsDTO := []*surveysv1.Question{}

	for _, question := range questions {
		questionDTO := surveysv1.Question{
			Id:       question.QuestionId.String(),
			Question: question.QuestionText,
			Type:     question.QuestionType,
		}

		questionsDTO = append(questionsDTO, &questionDTO)
	}
	resp := surveysv1.GetSurveyResp{
		Topic:    servey.Topic,
		SurveyId: servey.Id,
		Servey: &surveysv1.Survey{
			Question: questionsDTO,
		},
	}

	return &resp, nil
}

func (u ServeyUsecase) AddAnswers(ctx context.Context, in *surveysv1.AddAnswersReq) (*surveysv1.Nothing, error) {
	answers := in.Answer

	for _, answer := range answers {
		questionUUID, err := uuid.Parse(answer.QuestionId)
		if err != nil {
			return nil, err
		}
		userUUID, err := uuid.Parse(in.UserId)
		if err != nil {
			return nil, err
		}
		answerDTO := models.Answer{
			AnswerId:      uuid.New(),
			QuestionId:    questionUUID,
			UserId:        userUUID,
			TextAnswer:    *answer.TextAnswer,
			NumericAnswer: int(*answer.NumericAnswer),
		}

		err = u.repository.AddAnswer(ctx, answerDTO)
		if err != nil {
			return nil, err
		}
	}

	return &surveysv1.Nothing{
		Dummy: true,
	}, nil
}

const text = "text"
const numeric = "numeric"

func (u ServeyUsecase) GetStatictics(ctx context.Context, in *surveysv1.GetStaticticsReq) (*surveysv1.GetStaticticsResp, error) {

	questionId, err := uuid.Parse(in.Question_Id)
	if err != nil {
		return nil, err
	}

	questionType, err := u.repository.GetQuestionType(ctx, questionId)
	if err != nil {
		return nil, err
	}

	switch questionType {
	case text:
		questions, err := u.repository.GetAllTextAnswers(ctx, questionId)

		if err != nil {
			return nil, err
		}

		resp := surveysv1.GetStaticticsResp{
			QuestionsAnswerText: questions,
		}
		return &resp, nil
	case numeric:
		avg, err := u.repository.GetAverageNumericAnswer(ctx, questionId)

		if err != nil {
			return nil, err
		}

		resp := surveysv1.GetStaticticsResp{
			QuestionsAnswerNumeric: int64(avg),
		}
		return &resp, nil
	}

	return nil, fmt.Errorf("ало")
}
