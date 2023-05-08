package connect

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	MySQL = "mysql"
)

type Config struct {
	Driver   string
	User     string
	Password string
	DB       string
	Host     string
	Port     int
	Timeouts Timeouts
}

type Timeouts struct {
	Read time.Duration
}

func mysqlDSN(config Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Host, config.Port, config.DB)
}

func New(config Config) (db *sqlx.DB, err error) {
	switch config.Driver {
	case MySQL:
		db, err = sqlx.Open(MySQL, mysqlDSN(config))
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
