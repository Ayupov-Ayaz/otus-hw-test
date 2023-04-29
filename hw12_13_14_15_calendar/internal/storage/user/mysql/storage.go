package mysql

import (
	"context"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, user entity.User) (id int64, err error) {
	result, err := s.db.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)",
		user.Username,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *Storage) Get(ctx context.Context, id int64) (*entity.User, error) {
	user := &entity.User{}
	err := s.db.QueryRowxContext(ctx, "SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
