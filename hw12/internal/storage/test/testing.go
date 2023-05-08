package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

const (
	userName = "TEST_USER"
	userPass = "TEST_PWD"
	host     = "TEST_HOST"
	port     = "TEST_PORT"
)

func GetMysqlTestDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/calendar_test?parseTime=true",
		getUserName(), getPwd(), getHost(), getPort())
}

func getUserName() string {
	name := os.Getenv(userName)
	if name == "" {
		name = "user"
	}

	return name
}

func getPwd() string {
	return os.Getenv(userPass)
}

func getHost() string {
	host := os.Getenv(host)
	if host == "" {
		host = "localhost"
	}

	return host
}

func getPort() string {
	port := os.Getenv(port)
	if port == "" {
		port = "3306"
	}

	return port
}

// MysqlConnection - test connection to database.
func MysqlConnection(t *testing.T, dsn string) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Connect("mysql", dsn)
	require.NoError(t, err)

	err = db.PingContext(context.Background())
	require.NoError(t, err)

	return db
}
