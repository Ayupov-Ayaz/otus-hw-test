package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

type EventStorage interface {
	Create(ctx context.Context, event entity.Event) (id int64, err error)
	Update(ctx context.Context, event entity.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (entity.Event, error)
}

type Validator interface {
	Validate(ctx context.Context, i interface{}) error
}

type App struct {
	storage   EventStorage
	validator Validator
	logger    *zap.Logger
}

type Config func(*App)

func New(logger *zap.Logger, configs ...Config) *App {
	app := &App{
		logger: logger,
	}

	for _, cfg := range configs {
		cfg(app)
	}

	return app
}

func WithValidator(v Validator) Config {
	return func(a *App) {
		a.validator = v
	}
}

func WithStorage(s EventStorage) Config {
	return func(a *App) {
		a.storage = s
	}
}

func (a *App) CreateEvent(ctx context.Context, e entity.Event) error {
	if err := a.validator.Validate(ctx, e); err != nil {
		a.logger.Error("failed to validate event", zap.Error(err))
		return err
	}

	_, err := a.storage.Create(ctx, e)
	if err != nil {
		a.logger.Error("failed to create event", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event entity.Event) error {
	if err := a.storage.Update(ctx, event); err != nil {
		a.logger.Error("failed to update event",
			zap.Int64("id", event.ID),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	if err := a.storage.Delete(ctx, id); err != nil {
		a.logger.Error("failed to delete event",
			zap.Int64("id", id),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) GetEvent(ctx context.Context, id int64) (event entity.Event, err error) {
	event, err = a.storage.Get(ctx, id)
	if err != nil {
		a.logger.Error("failed to get event",
			zap.Int64("id", id),
			zap.Error(err))
	}

	return event, err
}
