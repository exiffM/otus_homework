package interfaces

import (
	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(format string, args ...any)
	Error(format string, args ...any)
	Warn(format string, args ...any)
	Debug(format string, args ...any)
}

type Storage interface {
	Connect() error
	Close() error
	CreateEvent(event mdl.Event) error
	SelectEvent(id int) (mdl.Event, error)
	UpdateEvent(event mdl.Event) error
	DeleteEvent(id int) error
	Events() ([]mdl.Event, error)
}

type Application interface { // TODO
}
