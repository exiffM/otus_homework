package memorystorage

import (
	"errors"
	"sort"
	"sync"

	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type LocalStorage = map[int]mdl.Event

type Storage struct {
	data LocalStorage
	mu   sync.RWMutex
}

var (
	errCreate = errors.New("event creation error")
	errExist  = errors.New("event doesn't exist")
)

func New() *Storage {
	return &Storage{nil, sync.RWMutex{}}
}

func (s *Storage) Connect() error {
	s.data = make(map[int]mdl.Event)
	return nil
}

func (s *Storage) Close() error {
	s.data = nil
	return nil
}

func (s *Storage) CreateEvent(event mdl.Event) error {
	if s.data == nil {
		return errCreate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	event.ID = len(s.data)
	s.data[event.ID] = event
	return nil
}

func (s *Storage) SelectEvent(id int) (mdl.Event, error) {
	if s.data == nil {
		return mdl.Event{}, errCreate
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.data[id]
	if !ok {
		return mdl.Event{}, errExist
	}
	return event, nil
}

func (s *Storage) UpdateEvent(event mdl.Event) error {
	if s.data == nil {
		return errCreate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[event.ID]
	if !ok {
		return errExist
	}
	s.data[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	if s.data == nil {
		return errCreate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
	return nil
}

func (s *Storage) Events() ([]mdl.Event, error) {
	if s.data == nil {
		return nil, errCreate
	}
	result := make([]mdl.Event, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		result = append(result, event)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
