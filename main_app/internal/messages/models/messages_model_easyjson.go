// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	uuid "github.com/google/uuid"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(in *jlexer.Lexer, out *WebScoketDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "messageType":
			out.MsgType = MsgType(in.String())
		case "payload":
			if m, ok := out.Payload.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Payload.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Payload = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(out *jwriter.Writer, in WebScoketDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messageType\":"
		out.RawString(prefix[1:])
		out.String(string(in.MsgType))
	}
	{
		const prefix string = ",\"payload\":"
		out.RawString(prefix)
		if m, ok := in.Payload.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Payload.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Payload))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WebScoketDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WebScoketDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WebScoketDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WebScoketDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(in *jlexer.Lexer, out *Payload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			out.URL = string(in.String())
		case "filename":
			out.Filename = string(in.String())
		case "size":
			out.Size = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(out *jwriter.Writer, in Payload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	{
		const prefix string = ",\"filename\":"
		out.RawString(prefix)
		out.String(string(in.Filename))
	}
	{
		const prefix string = ",\"size\":"
		out.RawString(prefix)
		out.Int64(int64(in.Size))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Payload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Payload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Payload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Payload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels1(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(in *jlexer.Lexer, out *MessagesArrayDTOOutput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "isNew":
			out.IsNew = bool(in.Bool())
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]Message, 0, 0)
					} else {
						out.Messages = []Message{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Message
					(v1).UnmarshalEasyJSON(in)
					out.Messages = append(out.Messages, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(out *jwriter.Writer, in MessagesArrayDTOOutput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"isNew\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.IsNew))
	}
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix)
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Messages {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesArrayDTOOutput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesArrayDTOOutput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesArrayDTOOutput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesArrayDTOOutput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels2(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(in *jlexer.Lexer, out *MessagesArrayDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]Message, 0, 0)
					} else {
						out.Messages = []Message{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Message
					(v4).UnmarshalEasyJSON(in)
					out.Messages = append(out.Messages, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(out *jwriter.Writer, in MessagesArrayDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix[1:])
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Messages {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesArrayDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesArrayDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesArrayDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesArrayDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels3(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(in *jlexer.Lexer, out *MessageInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(out *jwriter.Writer, in MessageInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels4(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(in *jlexer.Lexer, out *MessageDTOInput) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "disconnect":
			out.Disconnect = bool(in.Bool())
		case "message":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(out *jwriter.Writer, in MessageDTOInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"disconnect\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Disconnect))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageDTOInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageDTOInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageDTOInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageDTOInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels5(l, v)
}
func easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(in *jlexer.Lexer, out *Message) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "messageId":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.MessageId).UnmarshalText(data))
			}
		case "authorID":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.AuthorID).UnmarshalText(data))
			}
		case "branchId":
			if in.IsNull() {
				in.Skip()
				out.BranchID = nil
			} else {
				if out.BranchID == nil {
					out.BranchID = new(uuid.UUID)
				}
				if data := in.UnsafeBytes(); in.Ok() {
					in.AddError((*out.BranchID).UnmarshalText(data))
				}
			}
		case "text":
			out.Message = string(in.String())
		case "datetime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.SentAt).UnmarshalJSON(data))
			}
		case "chatId":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ChatId).UnmarshalText(data))
			}
		case "isRedacted":
			out.IsRedacted = bool(in.Bool())
		case "message_type":
			out.MessageType = string(in.String())
		case "parent_chat_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ChatIdParent).UnmarshalText(data))
			}
		case "files":
			if in.IsNull() {
				in.Skip()
				out.FilesDTO = nil
			} else {
				in.Delim('[')
				if out.FilesDTO == nil {
					if !in.IsDelim(']') {
						out.FilesDTO = make([]Payload, 0, 1)
					} else {
						out.FilesDTO = []Payload{}
					}
				} else {
					out.FilesDTO = (out.FilesDTO)[:0]
				}
				for !in.IsDelim(']') {
					var v7 Payload
					(v7).UnmarshalEasyJSON(in)
					out.FilesDTO = append(out.FilesDTO, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "photos":
			if in.IsNull() {
				in.Skip()
				out.PhotosDTO = nil
			} else {
				in.Delim('[')
				if out.PhotosDTO == nil {
					if !in.IsDelim(']') {
						out.PhotosDTO = make([]Payload, 0, 1)
					} else {
						out.PhotosDTO = []Payload{}
					}
				} else {
					out.PhotosDTO = (out.PhotosDTO)[:0]
				}
				for !in.IsDelim(']') {
					var v8 Payload
					(v8).UnmarshalEasyJSON(in)
					out.PhotosDTO = append(out.PhotosDTO, v8)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "sticker":
			out.Sticker = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messageId\":"
		out.RawString(prefix[1:])
		out.RawText((in.MessageId).MarshalText())
	}
	{
		const prefix string = ",\"authorID\":"
		out.RawString(prefix)
		out.RawText((in.AuthorID).MarshalText())
	}
	{
		const prefix string = ",\"branchId\":"
		out.RawString(prefix)
		if in.BranchID == nil {
			out.RawString("null")
		} else {
			out.RawText((*in.BranchID).MarshalText())
		}
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"datetime\":"
		out.RawString(prefix)
		out.Raw((in.SentAt).MarshalJSON())
	}
	{
		const prefix string = ",\"chatId\":"
		out.RawString(prefix)
		out.RawText((in.ChatId).MarshalText())
	}
	{
		const prefix string = ",\"isRedacted\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsRedacted))
	}
	{
		const prefix string = ",\"message_type\":"
		out.RawString(prefix)
		out.String(string(in.MessageType))
	}
	{
		const prefix string = ",\"parent_chat_id\":"
		out.RawString(prefix)
		out.RawText((in.ChatIdParent).MarshalText())
	}
	{
		const prefix string = ",\"files\":"
		out.RawString(prefix)
		if in.FilesDTO == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.FilesDTO {
				if v9 > 0 {
					out.RawByte(',')
				}
				(v10).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"photos\":"
		out.RawString(prefix)
		if in.PhotosDTO == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.PhotosDTO {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"sticker\":"
		out.RawString(prefix)
		out.String(string(in.Sticker))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2989b686EncodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2989b686DecodeGithubComGoParkMailRu20242EaglesDesignerMainAppInternalMessagesModels6(l, v)
}
