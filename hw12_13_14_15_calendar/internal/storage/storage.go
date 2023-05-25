package storage

import (
	"context"
	"io"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

type Event interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	GetEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error)
}

type Storage struct {
	event  Event
	closer io.Closer
}

func NewStorage(event Event, closer io.Closer) *Storage {
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
