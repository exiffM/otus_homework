package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	DEBUG   string = "debug"
	INFO    string = "info"
	WARNING string = "warn"
	ERROR   string = "error"
)

type Logger struct {
	level  string
	output io.Writer
}

func New(level string, out io.Writer) *Logger {
	return &Logger{level, out}
}

func (l *Logger) write(format string, args ...any) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[%s] %s ", l.level, time.Now().UTC().Format("01-12-2023 11:11")))
	sb.WriteString(fmt.Sprintf(format, args...))
	sb.WriteString("\n")
	l.output.Write([]byte(sb.String()))
}

func (l *Logger) Info(format string, args ...any) {
	if l.level == INFO {
		l.write(format, args...)
	}
}

func (l *Logger) Error(format string, args ...any) {
	if l.level == ERROR {
		l.write(format, args...)
	}
}

func (l *Logger) Warn(format string, args ...any) {
	if l.level == WARNING {
		l.write(format, args...)
	}
}

func (l *Logger) Debug(format string, args ...any) {
	if l.level == DEBUG {
		l.write(format, args...)
	}
}
