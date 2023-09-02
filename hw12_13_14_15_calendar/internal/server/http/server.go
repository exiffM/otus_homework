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

func NewServer(h, p string, l interfaces.Logger, a interfaces.Application) *Server {
	mux := mux.NewRouter().StrictSlash(true)
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(handleTeapot), l))
	mux.Handle("/api/calendar/event", api.NewCreateHandler(l, a))              // POST
	mux.Handle("/api/calendar/events/select/{id}", api.NewSelectHandler(l, a)) // GET
	mux.Handle("/api/calendar/events/update/{id}", api.NewUpdateHandler(l, a)) // PUT
	mux.Handle("/api/calendar/events/delete/{id}", api.NewDeleteHandler(l, a)) // DELETE
	mux.Handle("/api/calendar/events", api.NewEventsHandler(l, a))             // GET
	s := http.Server{
		Addr:              net.JoinHostPort(h, p),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           mux,
	}
	return &Server{&s, l, a}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		curCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := s.Stop(curCtx); err != nil {
			s.logger.Info("Server shutdown")
		}
	}()
	s.logger.Info("Server.Start()")
	err := http.ListenAndServe(s.server.Addr, s.server.Handler) //nolint: gosec
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
	s.logger.Info("Server.Stop()")
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
