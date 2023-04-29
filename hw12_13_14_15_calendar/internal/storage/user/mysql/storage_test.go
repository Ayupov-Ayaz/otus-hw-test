//go:build integration
// +build integration

package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/test"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/user/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

var testDSN = test.GetMysqlTestDSN()

func getConnection(t *testing.T) *sqlx.DB {
	return test.MysqlConnection(t, testDSN)
}

func remove(t *testing.T, db *sqlx.DB, id int64) {
	res, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	require.NoError(t, err)
	require.NotNil(t, res)

	rows, err := res.RowsAffected()
	require.Equal(t, int64(1), rows)
}

func create(t *testing.T, db *sqlx.DB, userName string) int64 {
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", userName)
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)

	return id
}

func getLastIDIfUserExist(t *testing.T, db *sqlx.DB) int64 {
	var lastID int64
	err := db.QueryRow("SELECT id FROM users ORDER BY id DESC LIMIT 1").Scan(&lastID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}

	return lastID
}

func removeUser(t *testing.T, db *sqlx.DB, id int64) {
	res, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	require.NoError(t, err)

	rows, err := res.RowsAffected()
	require.NoError(t, err)

	require.Equal(t, int64(1), rows)
}

func TestUserStorage_Create(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	userName := strconv.Itoa(time.Now().Nanosecond())
	storage := mysql.New(db)
	id, err := storage.Create(context.Background(), entity.User{ID: -1, Username: userName})
	require.NoError(t, err)
	require.True(t, id > 0)

	removeUser(t, db, id)
}

func TestUserStorage_Get(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	userName := strconv.Itoa(time.Now().Nanosecond())
	id := getLastIDIfUserExist(t, db)
	if id == 0 {
		id = create(t, db, userName)
	}

	storage := mysql.New(db)
	user, err := storage.Get(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, userName, user.Username)

	remove(t, db, id)
}
