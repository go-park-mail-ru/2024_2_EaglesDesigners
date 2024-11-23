package usecase

import (
	"context"

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
}

type ServeyUsecase struct {
	repository ServeyRepository
}

func (u *ServeyUsecase) GetSurvey(ctx context.Context, in *surveysv1.GetSurveyReq) (*surveysv1.GetSurveyResp, error) {
	serveyName := in.Name
	servey, err := u.repository.GetServey(ctx, serveyName)
	if err != nil {
		return nil, err
	}

	questions, err := u.repository.GetQuestionsByServeyName(ctx, serveyName)

	questionsDTO := []surveysv1.Question{}

	for _, question := range questions {
		questionDTO := surveysv1.Question{
			Id: int64(question.QuestionId[])
			SurveyId:
			Question:
			Type:
		}

		questionsDTO = append(questionsDTO, )
	}
	resp := surveysv1.GetSurveyResp{
		Topic: servey.Topic,
		SurveyId: servey.Id,
		Servey: &surveysv1.Survey{
			Question: []*surveysv1.Question{
				Id:
				SurveyId:
				Question:
				Type:
			},
		},
	}

	return &resp, nil
}

func (u *ServeyUsecase) AddAnswers(ctx context.Context, in *surveysv1.AddAnswersReq) (*surveysv1.Nothing, error) {

}
