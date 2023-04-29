package internal

import (
	"fmt"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
	memoryStorage "github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/event/memory"
	mysqlStorage "github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/event/mysql"
	"github.com/jmoiron/sqlx"
)

func NewEventStorage(driver string, db *sqlx.DB) (storage.Event, error) {
	var resp storage.Event

	switch driver {
	case MySQL:
		resp = mysqlStorage.New(db)
	case Memory:
		resp = memoryStorage.New()
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}

	return resp, nil
}
