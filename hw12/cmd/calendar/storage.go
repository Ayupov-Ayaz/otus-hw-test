package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/cmd/calendar/internal"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
)

func getConnect(config internal.StorageConf) (resp func() *sqlx.DB, err error) {
	switch config.Driver {
	case internal.Memory:
		resp = func() *sqlx.DB { return nil }
	default:
		var db *sqlx.DB
		db, err = internal.ConnectToDB(config)
		if err != nil {
			return nil, err
		}

		resp = func() *sqlx.DB { return db }
	}

	return resp, nil
}

func NewStorage(config internal.StorageConf) (*storage.Storage, error) {
	connection, err := getConnect(config)
	if err != nil {
		return nil, err
	}

	event, err := internal.NewEventStorage(config.Driver, connection())
	if err != nil {
		return nil, fmt.Errorf("failed to create event storage: %w", err)
	}

	return storage.NewStorage(event, connection()), nil
}
