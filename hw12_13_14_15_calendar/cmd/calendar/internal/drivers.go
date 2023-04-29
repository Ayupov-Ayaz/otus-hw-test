package internal

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	MySQL  = "mysql"
	Memory = "memory"
)

func MysqlDSN(config StorageConf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Host, config.Port, config.DB)
}

func ConnectToDB(config StorageConf) (db *sqlx.DB, err error) {
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
