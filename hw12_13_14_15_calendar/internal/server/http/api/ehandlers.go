package api

import (
	"net/http"
	"strconv"
	"strings"

	iface "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
	api "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http/api/default"
	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

type CreateHanlder struct {
	logger iface.Logger
	app    iface.Application
}

func NewCreateHandler(l iface.Logger, a iface.Application) *CreateHanlder {
	return &CreateHanlder{l, a}
}

func (ch *CreateHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := "api.create.event"
	if r.Method != "POST" {
		errString := "Invalid http method on create event"
		ch.logger.Info(errString)
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, errString),
			Logger:   ch.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
	var event mdl.Event
	if err := easyjson.UnmarshalFromReader(r.Body, &event); err != nil {
		ch.logger.Info("Unmarshal error on create event")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   ch.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	e, err := ch.app.CreateEvent(r.Context(), event)
	if err != nil {
		ch.logger.Info("Create event error")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   ch.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	response := api.NewEventResponse(method, "", e)
	_, _, err = easyjson.MarshalToHTTPResponseWriter(response, w)
	if err != nil {
		ch.logger.Info("Write http response error on create")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   ch.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
}

type SelectHanlder struct {
	logger iface.Logger
	app    iface.Application
}

func NewSelectHandler(l iface.Logger, a iface.Application) *SelectHanlder {
	return &SelectHanlder{l, a}
}

func (sh *SelectHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := "api.select.event"
	if r.Method != "GET" {
		errString := "Invalid http method on select event"
		sh.logger.Info(errString)
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, errString),
			Logger:   sh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sh.logger.Info("Id conversion error on select event")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   sh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	event, err := sh.app.SelectEvent(r.Context(), id)
	if err != nil {
		sh.logger.Info("Select event error")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   sh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	response := api.NewEventResponse(method, "", event)
	_, _, err = easyjson.MarshalToHTTPResponseWriter(response, w)
	if err != nil {
		sh.logger.Info("Write http response error on select")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   sh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
}

type UpdateHanlder struct {
	logger iface.Logger
	app    iface.Application
}

func NewUpdateHandler(l iface.Logger, a iface.Application) *UpdateHanlder {
	return &UpdateHanlder{l, a}
}

func (uh *UpdateHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := "api.update.event"
	if r.Method != "PUT" {
		errString := "Invalid http method on update event"
		uh.logger.Info(errString)
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, errString),
			Logger:   uh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	var event mdl.Event
	if err := easyjson.UnmarshalFromReader(r.Body, &event); err != nil {
		uh.logger.Info("Unmarshal error on update event")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   uh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		uh.logger.Info("Id conversion error on update event")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   uh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	event.ID = id
	_, err = uh.app.UpdateEvent(r.Context(), event)
	if err != nil {
		uh.logger.Info("Update event error")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   uh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	response := api.NewEventResponse(method, "", event)
	_, _, err = easyjson.MarshalToHTTPResponseWriter(response, w)
	if err != nil {
		uh.logger.Info("Write http response error on update")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   uh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
}

type DeleteHanlder struct {
	logger iface.Logger
	app    iface.Application
}

func NewDeleteHandler(l iface.Logger, a iface.Application) *DeleteHanlder {
	return &DeleteHanlder{l, a}
}

func (dh *DeleteHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := "api.delete.event"
	if r.Method != "DELETE" {
		errString := "Invalid http method on delete event"
		dh.logger.Info(errString)
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, errString),
			Logger:   dh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		dh.logger.Info("Id conversion error on delete event")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   dh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	err = dh.app.DeleteEvent(r.Context(), id)
	if err != nil {
		dh.logger.Info("Delete event error")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   dh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	sb := strings.Builder{}
	sb.WriteString("You've successfully deleted event with id: ")
	sb.WriteString(vars["id"])
	response := api.NewEventResponse(method, sb.String(), mdl.Event{})
	_, _, err = easyjson.MarshalToHTTPResponseWriter(response, w)
	if err != nil {
		dh.logger.Info("Write http response error on delete")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   dh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
}

type EventsHandler struct {
	logger iface.Logger
	app    iface.Application
}

func NewEventsHandler(l iface.Logger, a iface.Application) *EventsHandler {
	return &EventsHandler{l, a}
}

func (eh *EventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := "api.events.event"
	if r.Method != "GET" {
		errString := "Invalid http method on list events"
		eh.logger.Info(errString)
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, errString),
			Logger:   eh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	events, err := eh.app.Events()
	if err != nil {
		eh.logger.Info("List events error")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   eh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}

	response := api.NewEventsResponse(method, "", events)
	_, _, err = easyjson.MarshalToHTTPResponseWriter(response, w)
	if err != nil {
		eh.logger.Info("Write http response error on list events")
		handler := api.DefaultHandler{
			Response: *api.NewDefaultResponse(method, err.Error()),
			Logger:   eh.logger,
		}
		handler.ServeHTTP(w, r)
		return
	}
}
