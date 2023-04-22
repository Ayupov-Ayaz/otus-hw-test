package mysql

import (
	"context"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (s *EventRepository) Create(ctx context.Context, event entity.Event) (id int, err error) {
	// TODO: implement me
	panic("implement me")
}

func (s *EventRepository) Update(ctx context.Context, id int, event entity.Event) error {
	// TODO: implement me
	panic("implement me")

}

func (s *EventRepository) Delete(_ context.Context, id int) error {
	// TODO: implement me
	panic("implement me")

}

func (s *EventRepository) Get(ctx context.Context, id int) (entity.Event, error) {
	// TODO: implement me
	panic("implement me")

}
