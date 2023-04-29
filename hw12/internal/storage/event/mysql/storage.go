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
	createQuery = fmt.Sprintf("INSERT INTO events (%s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?)",
		titleField, userIDField, descriptionField, timeField, durationField, noticeField)
	updateQuery = fmt.Sprintf("UPDATE events SET %s = ?, %s = ?, %s = ?, %s = ?, %s = ? WHERE id = ?",
		titleField, descriptionField, timeField, durationField, noticeField)
	getQuery = fmt.Sprintf("SELECT id, %s, %s, %s, %s, %s, %s FROM events WHERE id = ?",
		titleField, userIDField, descriptionField, timeField, durationField, noticeField)
	deleteQuery = "DELETE FROM events WHERE id = ?"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) Create(ctx context.Context, event entity.Event) (id int64, err error) {
	result, err := s.db.ExecContext(ctx, createQuery, event.Title, event.UserID, event.Description, event.Time,
		event.DurationInSeconds(), event.BeforeStartNoticeInSeconds())

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

func (s *Repository) Update(ctx context.Context, event entity.Event) error {
	result, err := s.db.ExecContext(ctx, updateQuery, event.Title, event.Description, event.Time,
		event.DurationInSeconds(), event.BeforeStartNoticeInSeconds(), event.ID)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return checkExecResult(result)
}

func (s *Repository) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return checkExecResult(result)
}

func (s *Repository) Get(ctx context.Context, id int64) (entity.Event, error) {
	event := entity.Event{}
	var (
		duration int
		notice   int
	)

	err := s.db.QueryRowxContext(ctx, getQuery, id).
		Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &event.Time,
			&duration, &notice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrEventNotFound
		}

		return entity.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	event.Duration = time.Duration(duration) * time.Second
	event.BeforeStartNotice = time.Duration(notice) * time.Second

	return event, nil

}
