package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
	"github.com/jmoiron/sqlx"
)

const (
	titleField       = "title"
	descriptionField = "description"
	userIDField      = "user_id"
	timeField        = "time"
	durationField    = "duration_sec"
	noticeField      = "before_start_notice_sec"
)

var (
	createQuery = fmt.Sprintf("INSERT INTO events (%s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?)",
		titleField, userIDField, descriptionField, timeField, durationField)
	updateQuery = fmt.Sprintf("UPDATE events SET %s = ?, %s = ?, %s = ?, %s = ? WHERE id = ?",
		titleField, descriptionField, timeField, durationField)
	getQuery = fmt.Sprintf("SELECT id, %s, %s, %s, %s, %s FROM events WHERE id = ?",
		titleField, userIDField, descriptionField, timeField, durationField)
	deleteQuery = "DELETE FROM events WHERE id = ?"
)

// todo: СписокСобытийНаНеделю (дата начала недели)
// todo: СписокСобытийНaМесяц (дата начала месяца)
// todo: СписокСобытийНаДень (дата)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, event entity.Event) (id int64, err error) {
	result, err := s.db.ExecContext(ctx, createQuery, event.Title, event.UserID, event.Description, event.Time.Time(),
		event.DurationInSeconds())
	if err != nil {
		return 0, fmt.Errorf("failed to create event: %w", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
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

func (s *Storage) Update(ctx context.Context, event entity.Event) error {
	result, err := s.db.ExecContext(ctx, updateQuery, event.Title, event.Description, event.Time.Time(),
		event.DurationInSeconds(), event.ID)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return checkExecResult(result)
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, deleteQuery, id)
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

	err := s.db.QueryRowxContext(ctx, getQuery, id).
		Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &dateTime, &duration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrEventNotFound
		}

		return entity.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	event.Duration = entity.NewSecondsDuration(duration)
	event.Time = entity.MyTime(dateTime)

	return event, nil
}
