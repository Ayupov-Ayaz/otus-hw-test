package storage

import (
	"context"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
)

type Event interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	GetEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error)
}
