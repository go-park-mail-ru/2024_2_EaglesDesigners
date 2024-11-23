package api

import "context"

type Surveys interface {
	GetStatictics(ctx context.Context, in *surveysv1.GetStaticticsReq) (*surveysv1.GetStaticticsResp, error)
	GetSurvey(ctx context.Context, in *surveysv1.GetSurveyReq) (*surveysv1.GetSurveyResp, error)
	AddAnswers(ctx context.Context, in *surveysv1.AddAnswersReq) (*surveysv1.AddAnswersResp, error)
}

type Server struct {
	surveysv1.UnimplementedAuthServer
	surveys Surveys
}

func New(surveys Surveys) surveysv1.SurveysServer {
	return Server{
		auth: auth,
	}
}
