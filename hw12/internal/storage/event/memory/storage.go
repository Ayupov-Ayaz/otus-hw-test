package memory

import (
	"context"
	"sync"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int64]entity.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[int64]entity.Event),
		mu:     sync.RWMutex{},
	}
}

func (s *Storage) getEvent(id int64) (entity.Event, bool) {
	event, exist := s.events[id]
	return event, exist
}

func (s *Storage) Create(_ context.Context, event entity.Event) (id int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id = int64(len(s.events) + 1)
	s.events[id] = event

	return id, nil
}

func (s *Storage) Update(_ context.Context, event entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.getEvent(event.ID); !ok {
		return storage.ErrEventNotFound
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Get(_ context.Context, id int64) (entity.Event, error) {
	s.mu.RLock()
	event, ok := s.getEvent(id)
	s.mu.RUnlock()

	if !ok {
		return event, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.getEvent(id); !ok {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)

	return nil
}
