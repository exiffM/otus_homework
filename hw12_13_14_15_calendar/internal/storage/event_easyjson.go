// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package storage

import (
	json "encoding/json"
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

func easyjsonF642ad3eDecodeHw12131415CalendarInternalStorage(in *jlexer.Lexer, out *Event) {
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
		case "ID":
			out.ID = int(in.Int())
		case "Title":
			out.Title = string(in.String())
		case "Start":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Start).UnmarshalJSON(data))
			}
		case "Duration":
			out.Duration = int(in.Int())
		case "Description":
			out.Description = string(in.String())
		case "NotificationTime":
			out.NotificationTime = int(in.Int())
		case "Scheduled":
			out.Scheduled = bool(in.Bool())
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
func easyjsonF642ad3eEncodeHw12131415CalendarInternalStorage(out *jwriter.Writer, in Event) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Start\":"
		out.RawString(prefix)
		out.Raw((in.Start).MarshalJSON())
	}
	{
		const prefix string = ",\"Duration\":"
		out.RawString(prefix)
		out.Int(int(in.Duration))
	}
	{
		const prefix string = ",\"Description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"NotificationTime\":"
		out.RawString(prefix)
		out.Int(int(in.NotificationTime))
	}
	{
		const prefix string = ",\"Scheduled\":"
		out.RawString(prefix)
		out.Bool(bool(in.Scheduled))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Event) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF642ad3eEncodeHw12131415CalendarInternalStorage(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF642ad3eEncodeHw12131415CalendarInternalStorage(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeHw12131415CalendarInternalStorage(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeHw12131415CalendarInternalStorage(l, v)
}
