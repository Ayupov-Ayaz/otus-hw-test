package memory

import (
	"context"
	"sync"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
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
	if _, ok := s.getEvent(event.ID); ok {
		s.events[event.ID] = event
	}
	s.mu.Unlock()

	return nil
}

func (s *Storage) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetEventsForDates(_ context.Context, userID int64, start, end time.Time) ([]entity.Event, error) {
	var resp []entity.Event
	s.mu.RLock()
	for _, event := range s.events {
		eventDate := event.EventDate().YearDay()

		if event.UserID == userID &&
			eventDate >= start.YearDay() &&
			eventDate <= end.YearDay() {
			resp = append(resp, event)
		}
	}
	s.mu.RUnlock()

	return resp, nil
}
