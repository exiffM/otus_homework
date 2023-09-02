package api

import mdl "hw12_13_14_15_calendar/internal/storage"

//easyjson:json
type DefaultResponse struct {
	Method string
	Error  string
}

//easyjson:json
type EventResponse struct {
	Method string
	Error  string
	Data   mdl.Event
}

//easyjson:json
type EventsResponse struct {
	Method string
	Error  string
	Data   []mdl.Event
}

func NewDefaultResponse(m, e string) *DefaultResponse {
	return &DefaultResponse{Method: m, Error: e}
}

func NewEventResponse(m string, e string, d mdl.Event) *EventResponse {
	return &EventResponse{Method: m, Error: e, Data: d}
}

func NewEventsResponse(m string, e string, d []mdl.Event) *EventsResponse {
	return &EventsResponse{Method: m, Error: e, Data: d}
}
