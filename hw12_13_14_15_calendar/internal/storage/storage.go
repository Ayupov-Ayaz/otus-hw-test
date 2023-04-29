package storage

import (
	"context"
	"errors"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
)

var ErrEventNotFound = errors.New("event not found")

type Event interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (entity.Event, error)
}

type User interface {
	Create(ctx context.Context, user entity.User) (id int64, err error)
	Get(ctx context.Context, id int64) (*entity.User, error)
}

type Closer interface {
	Close() error
}

type Storage struct {
	event  Event
	user   User
	closer Closer
}

func NewStorage(event Event, user User, closer Closer) *Storage {
	return &Storage{
		event:  event,
		user:   user,
		closer: closer,
	}
}

func (s *Storage) Event() Event {
	return s.event
}

func (s *Storage) User() User {
	return s.user
}

func (s *Storage) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}

	return nil
}
