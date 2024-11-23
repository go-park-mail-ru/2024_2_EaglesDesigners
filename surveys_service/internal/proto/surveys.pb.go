// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: surveys_service/internal/proto/surveys.proto

package surveysv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Nothing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dummy bool `protobuf:"varint,1,opt,name=dummy,proto3" json:"dummy,omitempty"`
}

func (x *Nothing) Reset() {
	*x = Nothing{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Nothing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nothing) ProtoMessage() {}

func (x *Nothing) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nothing.ProtoReflect.Descriptor instead.
func (*Nothing) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{0}
}

func (x *Nothing) GetDummy() bool {
	if x != nil {
		return x.Dummy
	}
	return false
}

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionId    int64   `protobuf:"varint,2,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	TextAnswer    *string `protobuf:"bytes,3,opt,name=text_answer,json=textAnswer,proto3,oneof" json:"text_answer,omitempty"`
	NumericAnswer *int64  `protobuf:"varint,4,opt,name=numeric_answer,json=numericAnswer,proto3,oneof" json:"numeric_answer,omitempty"`
}

func (x *Answer) Reset() {
	*x = Answer{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{1}
}

func (x *Answer) GetQuestionId() int64 {
	if x != nil {
		return x.QuestionId
	}
	return 0
}

func (x *Answer) GetTextAnswer() string {
	if x != nil && x.TextAnswer != nil {
		return *x.TextAnswer
	}
	return ""
}

func (x *Answer) GetNumericAnswer() int64 {
	if x != nil && x.NumericAnswer != nil {
		return *x.NumericAnswer
	}
	return 0
}

type Question struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SurveyId int64  `protobuf:"varint,2,opt,name=survey_id,json=surveyId,proto3" json:"survey_id,omitempty"`
	Question string `protobuf:"bytes,3,opt,name=question,proto3" json:"question,omitempty"`
	TypeId   int64  `protobuf:"varint,4,opt,name=type_id,json=typeId,proto3" json:"type_id,omitempty"`
}

func (x *Question) Reset() {
	*x = Question{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Question) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Question) ProtoMessage() {}

func (x *Question) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Question.ProtoReflect.Descriptor instead.
func (*Question) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{2}
}

func (x *Question) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Question) GetSurveyId() int64 {
	if x != nil {
		return x.SurveyId
	}
	return 0
}

func (x *Question) GetQuestion() string {
	if x != nil {
		return x.Question
	}
	return ""
}

func (x *Question) GetTypeId() int64 {
	if x != nil {
		return x.TypeId
	}
	return 0
}

type Survey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question []*Question `protobuf:"bytes,1,rep,name=question,proto3" json:"question,omitempty"`
}

func (x *Survey) Reset() {
	*x = Survey{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Survey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Survey) ProtoMessage() {}

func (x *Survey) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Survey.ProtoReflect.Descriptor instead.
func (*Survey) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{3}
}

func (x *Survey) GetQuestion() []*Question {
	if x != nil {
		return x.Question
	}
	return nil
}

//	message GetStaticticsResp{
//	    repeated
//	}
type GetSurveyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetSurveyReq) Reset() {
	*x = GetSurveyReq{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSurveyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSurveyReq) ProtoMessage() {}

func (x *GetSurveyReq) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSurveyReq.ProtoReflect.Descriptor instead.
func (*GetSurveyReq) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{4}
}

func (x *GetSurveyReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetSurveyResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Servey   *Survey `protobuf:"bytes,1,opt,name=servey,proto3" json:"servey,omitempty"`
	Topic    string  `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	SurveyId int64   `protobuf:"varint,3,opt,name=survey_id,json=surveyId,proto3" json:"survey_id,omitempty"`
}

func (x *GetSurveyResp) Reset() {
	*x = GetSurveyResp{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSurveyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSurveyResp) ProtoMessage() {}

func (x *GetSurveyResp) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSurveyResp.ProtoReflect.Descriptor instead.
func (*GetSurveyResp) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{5}
}

func (x *GetSurveyResp) GetServey() *Survey {
	if x != nil {
		return x.Servey
	}
	return nil
}

func (x *GetSurveyResp) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *GetSurveyResp) GetSurveyId() int64 {
	if x != nil {
		return x.SurveyId
	}
	return 0
}

type AddAnswersReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64     `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Answer []*Answer `protobuf:"bytes,2,rep,name=answer,proto3" json:"answer,omitempty"`
}

func (x *AddAnswersReq) Reset() {
	*x = AddAnswersReq{}
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddAnswersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddAnswersReq) ProtoMessage() {}

func (x *AddAnswersReq) ProtoReflect() protoreflect.Message {
	mi := &file_surveys_service_internal_proto_surveys_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddAnswersReq.ProtoReflect.Descriptor instead.
func (*AddAnswersReq) Descriptor() ([]byte, []int) {
	return file_surveys_service_internal_proto_surveys_proto_rawDescGZIP(), []int{6}
}

func (x *AddAnswersReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AddAnswersReq) GetAnswer() []*Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

var File_surveys_service_internal_proto_surveys_proto protoreflect.FileDescriptor

var file_surveys_service_internal_proto_surveys_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x22, 0x1f, 0x0a, 0x07, 0x4e, 0x6f, 0x74, 0x68, 0x69,
	0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x22, 0x9e, 0x01, 0x0a, 0x06, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0b, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x61, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x74, 0x65, 0x78,
	0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x6e, 0x75,
	0x6d, 0x65, 0x72, 0x69, 0x63, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x01, 0x52, 0x0d, 0x6e, 0x75, 0x6d, 0x65, 0x72, 0x69, 0x63, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x5f,
	0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x6e, 0x75, 0x6d, 0x65, 0x72,
	0x69, 0x63, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x6c, 0x0a, 0x08, 0x51, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17,
	0x0a, 0x07, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x74, 0x79, 0x70, 0x65, 0x49, 0x64, 0x22, 0x37, 0x0a, 0x06, 0x53, 0x75, 0x72, 0x76, 0x65,
	0x79, 0x12, 0x2d, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x51, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x22, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x71,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6b, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x27, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e,
	0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x49,
	0x64, 0x22, 0x51, 0x0a, 0x0d, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x06, 0x61,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x75,
	0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x06, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x32, 0x81, 0x01, 0x0a, 0x07, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73,
	0x12, 0x3c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x12, 0x15, 0x2e,
	0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65,
	0x79, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x38,
	0x0a, 0x0a, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x73,
	0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x73, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x4e,
	0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b, 0x6e, 0x6f, 0x6e, 0x72,
	0x65, 0x70, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x73, 0x2e, 0x76, 0x31, 0x3b, 0x73, 0x75,
	0x72, 0x76, 0x65, 0x79, 0x73, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_surveys_service_internal_proto_surveys_proto_rawDescOnce sync.Once
	file_surveys_service_internal_proto_surveys_proto_rawDescData = file_surveys_service_internal_proto_surveys_proto_rawDesc
)

func file_surveys_service_internal_proto_surveys_proto_rawDescGZIP() []byte {
	file_surveys_service_internal_proto_surveys_proto_rawDescOnce.Do(func() {
		file_surveys_service_internal_proto_surveys_proto_rawDescData = protoimpl.X.CompressGZIP(file_surveys_service_internal_proto_surveys_proto_rawDescData)
	})
	return file_surveys_service_internal_proto_surveys_proto_rawDescData
}

var file_surveys_service_internal_proto_surveys_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_surveys_service_internal_proto_surveys_proto_goTypes = []any{
	(*Nothing)(nil),       // 0: surveys.Nothing
	(*Answer)(nil),        // 1: surveys.Answer
	(*Question)(nil),      // 2: surveys.Question
	(*Survey)(nil),        // 3: surveys.Survey
	(*GetSurveyReq)(nil),  // 4: surveys.GetSurveyReq
	(*GetSurveyResp)(nil), // 5: surveys.GetSurveyResp
	(*AddAnswersReq)(nil), // 6: surveys.AddAnswersReq
}
var file_surveys_service_internal_proto_surveys_proto_depIdxs = []int32{
	2, // 0: surveys.Survey.question:type_name -> surveys.Question
	3, // 1: surveys.GetSurveyResp.servey:type_name -> surveys.Survey
	1, // 2: surveys.AddAnswersReq.answer:type_name -> surveys.Answer
	4, // 3: surveys.Surveys.GetSurvey:input_type -> surveys.GetSurveyReq
	6, // 4: surveys.Surveys.AddAnswers:input_type -> surveys.AddAnswersReq
	5, // 5: surveys.Surveys.GetSurvey:output_type -> surveys.GetSurveyResp
	0, // 6: surveys.Surveys.AddAnswers:output_type -> surveys.Nothing
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_surveys_service_internal_proto_surveys_proto_init() }
func file_surveys_service_internal_proto_surveys_proto_init() {
	if File_surveys_service_internal_proto_surveys_proto != nil {
		return
	}
	file_surveys_service_internal_proto_surveys_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_surveys_service_internal_proto_surveys_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_surveys_service_internal_proto_surveys_proto_goTypes,
		DependencyIndexes: file_surveys_service_internal_proto_surveys_proto_depIdxs,
		MessageInfos:      file_surveys_service_internal_proto_surveys_proto_msgTypes,
	}.Build()
	File_surveys_service_internal_proto_surveys_proto = out.File
	file_surveys_service_internal_proto_surveys_proto_rawDesc = nil
	file_surveys_service_internal_proto_surveys_proto_goTypes = nil
	file_surveys_service_internal_proto_surveys_proto_depIdxs = nil
}
