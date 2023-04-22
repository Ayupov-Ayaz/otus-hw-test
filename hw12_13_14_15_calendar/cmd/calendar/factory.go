package main

import (
	"context"
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/sql/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewStorage(config StorageConf) (*storage.Storage, error) {
	var (
		db  *sqlx.DB
		err error
	)

	if !config.IsMemoryStorage() {
		db, err = connectToDB(config)
		if err != nil {
			return nil, err
		}
	}

	event, err := newEventRepository(config.Driver, db)
	if err != nil {
		return nil, err
	}

	closer, err := newCloser(config.Driver, db)
	if err != nil {
		return nil, err
	}

	return storage.NewStorage(event, closer), nil
}

func connectToDB(config StorageConf) (db *sqlx.DB, err error) {
	switch config.Driver {
	case MySQL:
		db, err = sqlx.Open(MySQL, MysqlDSN(config))
	default:
		err = fmt.Errorf("unknown driver: %s", config.Driver)
	}

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeouts.Read)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func MysqlDSN(config StorageConf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Host, config.Port, config.DB)
}

func newCloser(driver string, db *sqlx.DB) (storage.Closer, error) {
	switch driver {
	case MySQL:
		return db, nil
	case Memory:
		return memory.NewCloser(), nil
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}
}

func newEventRepository(driver string, db *sqlx.DB) (storage.EventRepository, error) {
	var resp storage.EventRepository

	switch driver {
	case MySQL:
		resp = mysql.NewEventRepository(db)
	case Memory:
		resp = memory.NewEventRepository()
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}

	return resp, nil
}
