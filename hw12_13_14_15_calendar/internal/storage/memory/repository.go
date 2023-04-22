package memory

import (
	"context"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"sync"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage"
)

type EventRepository struct {
	counter int
	mu      *sync.RWMutex
	list    map[int]entity.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		list: make(map[int]entity.Event),
		mu:   &sync.RWMutex{},
	}
}

func (s *EventRepository) getEvent(id int) (entity.Event, bool) {
	event, exist := s.list[id]
	return event, exist
}

func (s *EventRepository) Create(_ context.Context, event entity.Event) (id int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++

	event.ID = id
	s.list[s.counter] = event

	return id, nil
}

func (s *EventRepository) Update(_ context.Context, id int, event entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.getEvent(id)
	if !ok {
		return storage.ErrEventNotFound
	}

	s.list[id] = event

	return nil
}

func (s *EventRepository) Get(_ context.Context, id int) (entity.Event, error) {
	s.mu.RLock()
	event, ok := s.getEvent(id)
	s.mu.RUnlock()

	if !ok {
		return event, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *EventRepository) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.getEvent(id)
	if !ok {
		return storage.ErrEventNotFound
	}

	delete(s.list, id)

	return nil
}

func (s *EventRepository) Ping(_ context.Context) error {
	return nil
}
