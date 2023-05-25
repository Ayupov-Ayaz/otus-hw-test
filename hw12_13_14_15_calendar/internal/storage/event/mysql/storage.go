package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/jmoiron/sqlx"
)

var (
	ErrEventNotFound = errors.New("event not found")
)
var ErrNothingToUpdate = errors.New("nothing to update")

var (
	createQuery = "INSERT INTO events " +
		"(title, description, time, duration_sec,  before_notice_sec, user_id)" +
		" VALUES (?, ?, ?, ?, ?, ?)"

	getQuery = "SELECT id, title, description, time, duration_sec, before_notice_sec" +
		" FROM events " +
		"WHERE user_id = ? " +
		"AND time BETWEEN ? AND ?"
	deleteQuery = "DELETE FROM events WHERE id = ?"
)

type Storage struct {
	db     *sqlx.DB
	logger *zap.Logger
	qb     QueryBuilder
}

func New(db *sqlx.DB) *Storage {
	return &Storage{
		db:     db,
		logger: zap.L().Named("mysql event storage"),
	}
}

func (s *Storage) Create(ctx context.Context, event entity.Event) (id int64, err error) {
	result, err := s.db.ExecContext(ctx, createQuery,
		event.Title, event.Description, mySQLTimeFormat(event.DateTime.Time()),
		event.Duration.DurationInSec(), event.Notification.DurationInSec(), event.UserID)
	if err != nil {
		return 0, fmt.Errorf("failed to create event: %w", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert userID: %w", err)
	}

	return id, nil
}

func checkExecResult(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event entity.Event) error {
	eventQuery, ok := s.qb.updateEventQuery(event)
	if !ok {
		return ErrNothingToUpdate
	}

	res, err := s.db.ExecContext(ctx, eventQuery)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return checkExecResult(res)
}

func (s *Storage) Delete(ctx context.Context, id int64) (err error) {
	result, err := s.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return checkExecResult(result)
}

func (s *Storage) GetEventsForDates(ctx context.Context, userID int64, start, end time.Time) ([]entity.Event, error) {
	rows, err := s.db.QueryxContext(ctx, getQuery, userID, start, end)
	if err != nil {
		return nil, err
	}

	var events []entity.Event
	event := &entity.Event{}
	for rows.Next() {
		var (
			duration     int
			beforeNotice int
			dateTime     string
		)

		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &dateTime, &duration, &beforeNotice); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		event.Duration = entity.NewSecondsDuration(duration)
		date, err := parseMySQLTime(dateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %w", err)
		}

		event.DateTime = entity.MyTime(date)
		event.Notification = entity.NewSecondsDuration(beforeNotice)

		events = append(events, *event)
		event.Reset()
	}

	return events, nil
}
