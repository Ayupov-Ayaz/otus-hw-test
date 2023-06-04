package calendar

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
)

var ErrIDIsEmpty = errors.New("event id is empty")

type EventStorage interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	GetEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error)
}

type Validator interface {
	Validate(ctx context.Context, i interface{}) error
}

type EventUseCase struct {
	storage   EventStorage
	validator Validator
	logger    *zap.Logger
}

type Config func(*EventUseCase)

func New(validator Validator, storage EventStorage, logger *zap.Logger) *EventUseCase {
	return &EventUseCase{
		storage:   storage,
		validator: validator,
		logger:    logger,
	}
}

func (e *EventUseCase) CreateEvent(ctx context.Context, event entity.Event) (int64, error) {
	if err := e.validator.Validate(ctx, event); err != nil {
		e.logger.Error("failed to validate event", zap.Error(err))
		return 0, err
	}

	id, err := e.storage.Create(ctx, event)
	if err != nil {
		e.logger.Error("failed to create event", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (e *EventUseCase) UpdateEvent(ctx context.Context, event entity.Event) error {
	if err := e.validator.Validate(ctx, event); err != nil {
		e.logger.Error("failed to validate event", zap.Error(err))
		return err
	}

	if event.ID == 0 {
		e.logger.Error("failed to validate event: id is empty")
		return ErrIDIsEmpty
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := e.storage.Update(ctxWithTimeout, event); err != nil {
		e.logger.Error("failed to update event",
			zap.Int64("id", event.ID),
			zap.Error(err))
		return err
	}

	return nil
}

func (e *EventUseCase) DeleteEvent(ctx context.Context, id int64) error {
	if err := e.storage.Delete(ctx, id); err != nil {
		e.logger.Error("failed to delete event",
			zap.Int64("id", id),
			zap.Error(err))
		return err
	}

	return nil
}

func (e *EventUseCase) getEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error) {
	events, err := e.storage.GetEventsForDates(ctx, userID, start, end)
	if err != nil {
		e.logger.Error("failed to get events for dates",
			zap.Int64("userID", userID),
			zap.Time("start", start),
			zap.Time("end", end),
			zap.Error(err))
		return nil, err
	}

	return events, nil
}

func (e *EventUseCase) GetEventsByDay(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error) {
	return e.getEventsForDates(ctx, userID, date, date.AddDate(0, 0, 1))
}

func (e *EventUseCase) GetEventsByWeek(ctx context.Context, userID int64, startDate time.Time) ([]entity.Event, error) {
	return e.getEventsForDates(ctx, userID, startDate, startDate.AddDate(0, 0, 7))
}

func (e *EventUseCase) GetEventsByMonth(ctx context.Context, userID int64, startDate time.Time) ([]entity.Event, error) {
	return e.getEventsForDates(ctx, userID, startDate, startDate.AddDate(0, 1, 0))
}
