package models

// message Answer {
// 	int64 question_id = 2;
// 	optional string text_answer = 3;
// 	optional int64 numeric_answer = 4;
//   }

//   message Question {
// 	int64 id = 1;
// 	int64 survey_id = 2;
// 	string question = 3;
// 	string type = 4;
//   }

//   message Survey { repeated Question question = 1; }

//   // message GetStaticticsResp{
//   //     repeated
//   // }
//   message GetSurveyReq { string name = 1; }

//   message GetSurveyResp {
// 	Survey servey = 1;
// 	string topic = 2;
// 	int64 survey_id = 3;
//   }
//   message AddAnswersReq {
// 	int64 user_id = 1;
// 	repeated Answer answer = 2;
//   }

//   service Surveys {
// 	//   rpc GetStatictics(Nothing) returns {GetStaticticsResp}
// 	rpc GetSurvey(GetSurveyReq) returns (GetSurveyResp) {}
// 	rpc AddAnswers(AddAnswersReq) returns (Nothing) {}
//   }

type QuestionDTO struct {
	ID           string `json:"id" valid:"-"`
	Question     string `json:"question" valid:"-"`
	QuestionType string `json:"type" valid:"-"`
}

type GetSurveyDTO struct {
	Questions QuestionDTO `json:"questions" valid:"-"`
	Topic     string      `json:"topic" valid:"-"`
	Survey_id string      `json:"survey_id" valid:"-"`
}
