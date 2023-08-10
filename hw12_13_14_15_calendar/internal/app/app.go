package app

import (
	"context"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
)

type App struct { // TODO
}

func New(logger interfaces.Logger, storage interfaces.Storage) *App {
	_ = logger
	_ = storage
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	_ = ctx
	_ = id
	_ = title
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
