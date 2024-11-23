package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/surveys/models"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/utils/responser"
	surveyv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/protos/gen/go/surveyv1"
	"github.com/google/uuid"
	"go.octolab.org/pointer"

	auth "github.com/go-park-mail-ru/2024_2_EaglesDesigner/main_app/internal/auth/models"

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
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	vars := mux.Vars(r)
	questionId := vars["questionId"]

	grpcResp, err := d.client.GetStatictics(ctx, &surveyv1.GetStaticticsReq{
		Question_Id: questionId,
	})

	if err != nil {
		log.Errorf("не удалось получить статистику")
		responser.SendError(ctx, w, responser.InvalidJSONError, http.StatusBadRequest)
	}

	resp := convertFromGRPCStatictics(grpcResp)

	responser.SendStruct(ctx, w, resp, http.StatusOK)
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

	resp := convertFromGRPCSurvey(grpcResp)

	responser.SendStruct(ctx, w, resp, http.StatusOK)
}

type AnswersInput struct {
	Answers []Answer `json:"answers"`
}

type Answer struct {
	TextAnswer    string `json:"textAnswer"`
	NumericAnswer int64  `json:"numeric"`
}

// question/{questionId} post
func (d *Delivery) AddAnswers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	mapVars, ok := r.Context().Value(auth.MuxParamsKey).(map[string]string)
	if !ok {
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	questionId, err := uuid.Parse(mapVars["questionId"])

	var messageDTO AnswersInput
	err = json.NewDecoder(r.Body).Decode(&messageDTO)
	if err != nil {
		responser.SendError(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(auth.UserKey).(auth.User)
	log.Println(user)
	if !ok {
		log.Println("Message delivery -> AddNewMessage: нет юзера в контексте")
		responser.SendError(ctx, w, "Нет нужных параметров", http.StatusInternalServerError)
		return
	}

	answerArray := []*surveyv1.Answer{}

	for _, answer := range messageDTO.Answers {
		answerArray = append(answerArray,
			&surveyv1.Answer{
				QuestionId:    questionId.String(),
				TextAnswer:    &answer.TextAnswer,
				NumericAnswer: &answer.NumericAnswer,
			})
	}

	request := surveyv1.AddAnswersReq{
		UserId: user.ID.String(),
		Answer: answerArray,
	}

	grpcResp, err := d.client.AddAnswers(
		ctx,
		&request,
	)
	if err != nil {
		responser.SendError(ctx, w, "Не удалось добавить ответы", 500)
		return
	}

	res := grpcResp.Dummy

	if res {
		responser.SendOK(w, "все ок", 201)
		return
	}
	responser.SendError(ctx, w, "Не удалось добавить ответы", 500)
}

func convertFromGRPCSurvey(survey *surveyv1.GetSurveyResp) models.GetSurveyDTO {
	return models.GetSurveyDTO{
		Questions: convertFromGRPCQuestions(survey.GetServey().GetQuestion()),
		Topic:     survey.GetTopic(),
		Survey_id: survey.GetSurveyId(),
	}
}

func convertFromGRPCQuestions(questions []*surveyv1.Question) []models.QuestionDTO {
	var questionsDTO []models.QuestionDTO

	for _, question := range questions {
		questionDTO := models.QuestionDTO{
			ID:           question.GetId(),
			Question:     question.GetQuestion(),
			QuestionType: question.GetType(),
		}
		questionsDTO = append(questionsDTO, questionDTO)
	}

	return questionsDTO
}

func convertFromGRPCStatictics(statistics *surveyv1.GetStaticticsResp) models.GetStatictics {
	return models.GetStatictics{
		TextAnswers:    statistics.GetQuestionsAnswerText(),
		NumericAnswers: pointer.ToInt(int(statistics.GetQuestionsAnswerNumeric())),
	}
}
