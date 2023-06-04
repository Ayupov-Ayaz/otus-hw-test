package connect

import (
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/event"

	config "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/storage"
	"github.com/jmoiron/sqlx"
)

func getConnect(config config.Config) (resp func() *sqlx.DB, err error) {
	if config.Driver == "memory" {
		return func() *sqlx.DB { return nil }, nil
	}

	var db *sqlx.DB
	db, err = New(Config(config))
	if err != nil {
		return nil, err
	}

	return func() *sqlx.DB { return db }, nil
}

func NewStorage(config config.Config) (storage.Event, error) {
	connection, err := getConnect(config)
	if err != nil {
		return nil, err
	}

	eventStorage, err := event.New(config.Driver, connection())
	if err != nil {
		return nil, fmt.Errorf("failed to create event storage: %w", err)
	}

	return eventStorage, nil
}
