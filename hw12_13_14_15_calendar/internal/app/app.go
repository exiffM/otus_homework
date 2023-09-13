package app

import (
	"context"

	interfaces "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  interfaces.Logger
	storage interfaces.Storage
}

func New(l interfaces.Logger, s interfaces.Storage) *App {
	return &App{l, s}
}

func (a *App) CreateEvent(ctx context.Context, event mdl.Event) (mdl.Event, error) {
	// TODO
	_ = ctx
	if err := a.storage.Connect(); err != nil {
		return mdl.Event{}, err
	}
	defer a.storage.Close()
	return a.storage.CreateEvent(event)
}

func (a *App) SelectEvent(ctx context.Context, id int) (mdl.Event, error) {
	_ = ctx
	if err := a.storage.Connect(); err != nil {
		return mdl.Event{}, err
	}
	defer a.storage.Close()
	return a.storage.SelectEvent(id)
}

func (a *App) UpdateEvent(ctx context.Context, event mdl.Event) (mdl.Event, error) {
	_ = ctx
	if err := a.storage.Connect(); err != nil {
		return mdl.Event{}, err
	}
	defer a.storage.Close()
	return a.storage.UpdateEvent(event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	_ = ctx
	if err := a.storage.Connect(); err != nil {
		return err
	}
	defer a.storage.Close()
	return a.storage.DeleteEvent(id)
}

func (a *App) Events() ([]mdl.Event, error) {
	if err := a.storage.Connect(); err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.Events()
}

func (a *App) NotScheduledEvents() ([]mdl.Event, error) {
	if err := a.storage.Connect(); err != nil {
		return nil, err
	}
	defer a.storage.Close()
	return a.storage.NotScheduledEvents()
}
