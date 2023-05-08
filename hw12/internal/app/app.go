package app

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

var ErrIDIsEmpty = errors.New("event id is empty")

type EventStorage interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (entity.Event, error)
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

func New(logger *zap.Logger, configs ...Config) *EventUseCase {
	app := &EventUseCase{
		logger: logger,
	}

	for _, cfg := range configs {
		cfg(app)
	}

	return app
}

func WithValidator(v Validator) Config {
	return func(e *EventUseCase) {
		e.validator = v
	}
}

func WithStorage(s EventStorage) Config {
	return func(e *EventUseCase) {
		e.storage = s
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

func (e *EventUseCase) GetEvent(ctx context.Context, id int64) (event entity.Event, err error) {
	event, err = e.storage.Get(ctx, id)
	if err != nil {
		e.logger.Error("failed to get event",
			zap.Int64("id", id),
			zap.Error(err))
	}

	return event, err
}
