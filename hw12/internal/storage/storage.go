package storage

import (
	"context"
	"errors"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

var ErrEventNotFound = errors.New("event not found")

type Event interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (entity.Event, error)
}

type Closer interface {
	Close() error
}

type Storage struct {
	event  Event
	closer Closer
}

func NewStorage(event Event, closer Closer) *Storage {
	return &Storage{
		event:  event,
		closer: closer,
	}
}

func (s *Storage) Event() Event {
	return s.event
}

func (s *Storage) Close() error {
	if s.closer != nil {
		return s.closer.Close()
	}

	return nil
}
