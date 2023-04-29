package main

import (
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

func NewStorage(config internal.StorageConf) (*storage.Storage, error) {
	var (
		db  *sqlx.DB
		err error
	)

	if !config.IsMemoryStorage() {
		db, err = internal.ConnectToDB(config)
		if err != nil {
			return nil, err
		}
	}

	event, err := internal.NewEventStorage(config.Driver, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create event storage: %w", err)
	}

	user, err := internal.NewUserStorage(config.Driver, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create user storage: %w", err)
	}

	return storage.NewStorage(event, user, db), nil
}
