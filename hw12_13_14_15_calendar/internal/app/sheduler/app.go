package sheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(event string, body []byte) error
}

type Storage interface {
	GetEventsForDates(ctx context.Context, start, end time.Time) ([]entity.Event, error)
	Delete(ctx context.Context, id int64) error
}

type App struct {
	storage   Storage
	publisher Publisher
	logger    *zap.Logger
}

func New(storage Storage, publisher Publisher, logger *zap.Logger) *App {
	return &App{
		storage:   storage,
		publisher: publisher,
		logger:    logger,
	}
}

func (a *App) Sender(eventID string, addDuration time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		now := time.Now()
		events, err := a.storage.GetEventsForDates(ctx, now, now.Add(addDuration))
		if err != nil {
			return fmt.Errorf("failed to get events for dates: %w", err)
		}

		for _, event := range events {
			data, err := jsoniter.Marshal(event)
			if err != nil {
				return fmt.Errorf("failed to marshal event: %w", err)
			}

			if err := a.publisher.Publish(eventID, data); err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}

			if err := a.storage.Delete(ctx, event.ID); err != nil {
				return fmt.Errorf("failed to delete event with id=%d: %w", event.ID, err)
			}
		}

		return nil
	}
}

func (a *App) Start(done <-chan struct{}, name string, work func(ctx context.Context) error, interval time.Duration) {
	logger := a.logger.Named("worker_" + name)

	logger.Info("started", zap.String("interval", interval.String()))
	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			if err := work(context.Background()); err != nil {
				a.logger.Error("failed to work", zap.Error(err))
			}
		case <-done:
			a.logger.Info("worker stopped")
			return
		default:
			//
		}
	}
}

func (a *App) remove(ctx context.Context) error {
	startDate := time.Now().AddDate(-1, 0, 0)
	endDate := time.Now().AddDate(0, 0, -1)
	events, err := a.storage.GetEventsForDates(ctx, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get events for dates: %w", err)
	}

	for _, event := range events {
		if err := a.storage.Delete(ctx, event.ID); err != nil {
			return fmt.Errorf("failed to delete event with id=%d: %w", event.ID, err)
		}
	}

	return nil
}

func (a *App) Remover(done <-chan struct{}, interval time.Duration) {
	a.logger.Info("remover started", zap.String("interval", interval.String()))
	tick := time.Tick(interval)
	for {
		select {
		case <-tick:
			if err := a.remove(context.Background()); err != nil {
				a.logger.Error("failed to remove", zap.Error(err))
			}
		case <-done:
			a.logger.Info("remover stopped")
			return
		default:
			//
		}
	}
}
