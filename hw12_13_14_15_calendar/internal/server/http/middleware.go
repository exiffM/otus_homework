package internalhttp

import (
	"net/http"
	"time"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
)

type ResponceLogger struct {
	http.ResponseWriter
	code int
}

func loggingMiddleware(next http.Handler, logger interfaces.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startAt := time.Now()
		rl := NewLoggingResponseWriter(w)
		next.ServeHTTP(rl, r)
		a := struct {
			UserAgent       string
			ClientIPAddress string
			HTTPMethod      string
			HTTPVersion     string
			StartAt         time.Time
			Latency         time.Duration
		}{
			UserAgent:       r.UserAgent(),
			ClientIPAddress: r.RemoteAddr,
			HTTPMethod:      r.Method,
			HTTPVersion:     r.Proto,
			StartAt:         startAt,
			Latency:         time.Since(startAt),
		}
		logger.Info("%+v", a)
	})
}

func NewLoggingResponseWriter(writer http.ResponseWriter) *ResponceLogger {
	return &ResponceLogger{writer, 0}
}

func (rl *ResponceLogger) WriteHeader(code int) {
	rl.code = code
	rl.ResponseWriter.WriteHeader(code)
}
