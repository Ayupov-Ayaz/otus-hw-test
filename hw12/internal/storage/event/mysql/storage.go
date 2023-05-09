package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db     *sqlx.DB
	logger *zap.Logger
	qb     QueryBuilder
}

func New(db *sqlx.DB, logger *zap.Logger) *Storage {
	return &Storage{
		db:     db,
		logger: logger,
	}
}

func (s *Storage) rollbackTx(tx *sql.Tx) {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		s.logger.Error("failed to rollback transaction", zap.Error(rollbackErr))
	}
}

func (s *Storage) createNotifications(ctx context.Context, tx *sql.Tx, eventID int64, notifications []entity.Duration) error {
	notificationQuery := s.qb.createNotificationQuery(eventID, notifications)
	if _, err := tx.ExecContext(ctx, notificationQuery); err != nil {
		return fmt.Errorf("failed to create notifications: %w", err)
	}

	return nil
}

func (s *Storage) Create(ctx context.Context, event entity.Event) (id int64, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			s.rollbackTx(tx)
		}
	}()

	result, err := tx.ExecContext(ctx,
		"INSERT INTO events (title, description, time, duration_sec, user_id) VALUES (?, ?, ?, ?, ?)",
		event.Title, event.Description, mySQLTimeFormat(event.DateTime.Time()), event.Duration.DurationInSec(), event.UserID)
	if err != nil {
		return 0, fmt.Errorf("failed to create event: %w", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	if err := s.createNotifications(ctx, tx, id, event.Notifications); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func checkExecResult(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) deleteNotifications(ctx context.Context, tx *sql.Tx, eventID int64) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM notifications WHERE event_id = ?", eventID)
	if err != nil {
		return fmt.Errorf("delete eventID=%d notifications failed: %w", eventID, err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event entity.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			s.rollbackTx(tx)
		}
	}()

	eventQuery, ok := s.qb.updateEventQuery(event)
	if ok {
		_, err := tx.ExecContext(ctx, eventQuery)
		if err != nil {
			return fmt.Errorf("failed to update event: %w", err)
		}
	}

	if len(event.Notifications) > 0 {
		if err := s.deleteNotifications(ctx, tx, event.ID); err != nil {
			return err
		}

		if err := s.createNotifications(ctx, tx, event.ID, event.Notifications); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	defer func() {
		if err != nil {
			s.rollbackTx(tx)
		}
	}()

	if err := s.deleteNotifications(ctx, tx, id); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM events WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return checkExecResult(result)
}

func (s *Storage) Get(ctx context.Context, id int64) (entity.Event, error) {
	event := entity.Event{}
	var (
		duration int
		dateTime time.Time
	)

	err := s.db.QueryRowxContext(ctx, "SELECT FROM events `id`, `title`, `user_id`, `description`, `time`, "+
		"`duration_sec` WHERE id = ?", id).
		Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &dateTime, &duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrEventNotFound
		}

		return entity.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	event.Duration = entity.NewSecondsDuration(duration)
	event.DateTime = entity.MyTime(dateTime)

	return event, nil
}

func (s *Storage) GetEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error) {
	rows, err := s.db.QueryxContext(ctx, "SELECT id, title, description, time, duration_sec FROM events WHERE user_id = ? "+
		"AND time BETWEEN ? AND ?", userID, start, end)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrEventNotFound
		}
		return nil, err
	}

	var events []entity.Event
	event := &entity.Event{}
	for rows.Next() {
		var (
			duration int
			dateTime string
		)

		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &dateTime, &duration); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		event.Duration = entity.NewSecondsDuration(duration)
		date, err := parseMySQLTime(dateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %w", err)
		}

		event.DateTime = entity.MyTime(date)
		events = append(events, *event)
		event.Reset()
	}

	return events, nil
}
