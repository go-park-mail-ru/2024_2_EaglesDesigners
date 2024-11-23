package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	surveyv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/surveyv1"
)

type Delivery struct {
	client surveyv1.SurveysClient
}

func New(client surveyv1.SurveysClient) *Delivery {
	return &Delivery{
		client: client,
	}
}

func (d *Delivery) GetStatictics(w http.ResponseWriter, r *http.Request) {

}

func (d *Delivery) GetSurvey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	
}

func (d *Delivery) AddAnswers(w http.ResponseWriter, r *http.Request) {

}
