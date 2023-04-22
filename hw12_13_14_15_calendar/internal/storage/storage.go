package storage

import (
	"context"
	"errors"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventRepository interface {
	Create(ctx context.Context, event entity.Event) (id int, err error)
	Update(ctx context.Context, id int, event entity.Event) error
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (entity.Event, error)
}

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (id int, err error)
}

type Closer interface {
	Close() error
}

type Storage struct {
	event  EventRepository
	closer Closer
}

func NewStorage(event EventRepository, closer Closer) *Storage {
	return &Storage{
		event:  event,
		closer: closer,
	}
}

func (s *Storage) Event() EventRepository {
	return s.event
}

func (s *Storage) Close() error {
	return s.closer.Close()
}
