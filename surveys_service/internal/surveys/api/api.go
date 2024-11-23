package api

type Surveys interface {
	GetStatictics(ctx context.Context, in *surveysv1.) (*surveysv1. ,error)
	GetSurvey(ctx context.Context, in *surveysv1.) (*surveysv1. ,error) 
    AddAnswers(ctx context.Context, in *surveysv1.) (*surveysv1. ,error)
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



