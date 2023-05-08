package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/event"

	config "github.com/ayupov-ayaz/otus-wh-test/hw12/cmd/calendar/internal/configs/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/connect"
	"github.com/jmoiron/sqlx"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
)

func getConnect(config config.Config) (resp func() *sqlx.DB, err error) {
	if config.Driver == "memory" {
		return func() *sqlx.DB { return nil }, nil
	}

	cfg := connect.Config{
		Driver:   config.Driver,
		User:     config.User,
		Password: config.Password,
		DB:       config.DB,
		Host:     config.Host,
		Port:     config.Port,
		Timeouts: connect.Timeouts{Read: config.Timeouts.Read},
	}

	var db *sqlx.DB
	db, err = connect.New(cfg)
	if err != nil {
		return nil, err
	}

	return func() *sqlx.DB { return db }, nil
}

func NewStorage(config config.Config, logger *zap.Logger) (*storage.Storage, error) {
	connection, err := getConnect(config)
	if err != nil {
		return nil, err
	}

	eventStorage, err := event.New(config.Driver, connection(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create event storage: %w", err)
	}

	return storage.NewStorage(eventStorage, connection()), nil
}
