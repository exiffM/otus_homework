package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
	"github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/server/http/api"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	logger interfaces.Logger
	app    interfaces.Application
}

func NewServer(h, p string, rht int, l interfaces.Logger, a interfaces.Application) *Server {
	mux := mux.NewRouter().StrictSlash(true)
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(handleTeapot), l))
	mux.Handle("/api/calendar/event", api.NewCreateHandler(l, a))              // POST
	mux.Handle("/api/calendar/events/select/{id}", api.NewSelectHandler(l, a)) // GET
	mux.Handle("/api/calendar/events/update/{id}", api.NewUpdateHandler(l, a)) // PUT
	mux.Handle("/api/calendar/events/delete/{id}", api.NewDeleteHandler(l, a)) // DELETE
	mux.Handle("/api/calendar/events", api.NewEventsHandler(l, a))             // GET
	s := http.Server{
		Addr:              net.JoinHostPort(h, p),
		ReadHeaderTimeout: time.Duration(rht) * time.Second,
		ReadTimeout:       time.Duration(rht) * time.Second,
		Handler:           mux,
	}
	return &Server{&s, l, a}
}

func (s *Server) Start() error {
	s.logger.Info("Server.Start()")
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func handleTeapot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Status code has been received!"))
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP Server.Stop()")
	err := s.server.Shutdown(ctx)
	s.logger.Error("After HTTP shutdown!")
	if err != nil {
		return err
	}
	return nil
}
