package connect

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	MySQL   = "mysql"
	timeout = 5 * time.Second
)

type Config struct {
	Driver   string
	User     string
	Password string
	DB       string
	Host     string
	Port     int
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
