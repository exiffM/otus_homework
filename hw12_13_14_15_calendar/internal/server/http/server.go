package internalhttp

import (
	"context"
	"net/http"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
)

type Server struct {
	host   string
	port   string
	logger interfaces.Logger
	app    interfaces.Application
}

func NewServer(h, p string, l interfaces.Logger, a interfaces.Application) *Server {
	return &Server{h, p, l, a}
}

func (s *Server) Start(ctx context.Context) error {
	_ = ctx
	// <-ctx.Done()
	s.logger.Info("Server.Start()")
	http.Handle("/", loggingMiddleware(http.HandlerFunc(handleTeapot), s.logger))
	err := http.ListenAndServe(s.host+":"+s.port, nil) //nolint: gosec
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
	// TODO
	_ = ctx
	s.logger.Info("Server.Start()")
	return nil
}
