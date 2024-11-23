package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/responser"
	surveyv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/surveyv1"
	"github.com/gorilla/mux"
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

// @Router /survey/{name} [get]
func (d *Delivery) GetSurvey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	vars := mux.Vars(r)
	surveyName := vars["name"]

	grpcResp, err := d.client.GetSurvey(
		ctx,
		&surveyv1.GetSurveyReq{
			Name: surveyName,
		},
	)

	if err != nil {
		log.Errorf("не удалось получить опрос")
		responser.SendError(ctx, w, responser.InvalidJSONError, http.StatusBadRequest)
	}

	responser.SendStruct(ctx, w, resp, http.StatusOK)
}

func (d *Delivery) AddAnswers(w http.ResponseWriter, r *http.Request) {

}

func convertFromGRPCSurvey(survey surveyv1.GetSurveyResp) models.
