package api

import (
	"context"

	surveysv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/proto"
)

type Surveys interface {
	GetStatictics(ctx context.Context, in *surveysv1.GetStaticticsReq) (*surveysv1.GetStaticticsResp, error)
	GetSurvey(ctx context.Context, in *surveysv1.GetSurveyReq) (*surveysv1.GetSurveyResp, error)
	AddAnswers(ctx context.Context, in *surveysv1.AddAnswersReq) (*surveysv1.Nothing, error)
}

type Server struct {
	surveysv1.UnimplementedSurveysServer
	surveys Surveys
}

func New(surveys Surveys) surveysv1.SurveysServer {
	return Server{
		surveys: surveys,
	}
}

func (s Server) GetStatictics(ctx context.Context, in *surveysv1.GetStaticticsReq) (*surveysv1.GetStaticticsResp, error) {
	return s.surveys.GetStatictics(ctx, in)
}

func (s Server) GetSurvey(ctx context.Context, in *surveysv1.GetSurveyReq) (*surveysv1.GetSurveyResp, error) {
	return s.surveys.GetSurvey(ctx, in)
}
func (s Server) AddAnswers(ctx context.Context, in *surveysv1.AddAnswersReq) (*surveysv1.Nothing, error) {
	return s.surveys.AddAnswers(ctx, in)
}
