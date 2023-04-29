//go:build integration
// +build integration

package mysql

import (
	"context"
	store "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	testDSN  string
	dateTime time.Time
)

func TestMain(m *testing.M) {
	testDSN = test.GetMysqlTestDSN()
	var err error
	dateTime, err = time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func getConnection(t *testing.T) *sqlx.DB {
	return test.GetMysqlConnection(t, testDSN)
}

func createUser(t *testing.T, db *sqlx.DB) entity.User {
	userName := strconv.Itoa(time.Now().Nanosecond())
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", userName)
	require.NoError(t, err, "create user failed")
	id, err := res.LastInsertId()
	require.NoError(t, err, "get last insert id failed")

	return entity.NewUser(id, userName)
}

func deleteUserByUserName(t *testing.T, db *sqlx.DB, userName string) {
	res, err := db.Exec("DELETE FROM users WHERE name = ?", userName)
	require.NoError(t, err, "delete user failed")
	rows, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rows, "user not deleted")
}

func deleteEvent(t *testing.T, db *sqlx.DB, id int64) {
	res, err := db.Exec(deleteQuery, id)
	require.NoError(t, err)
	rows, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
}

func createEvent(t *testing.T, db *sqlx.DB, e entity.Event) int64 {
	res, err := db.Exec(createQuery, e.Title, e.UserID, e.Description, e.Time,
		e.DurationInSeconds(), e.BeforeStartNoticeInSeconds())
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return id
}

func makeEvent(userID int64) entity.Event {
	duration := 5 * time.Second
	beforeDuration := 10 * time.Second
	title := strconv.Itoa(time.Now().Nanosecond())
	return entity.NewEvent(title, "desc", userID, dateTime, duration, beforeDuration)
}

func TestEventRepository_Create(t *testing.T) {
	makeEvent(1)
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	user := createUser(t, db)
	defer func() {
		deleteUserByUserName(t, db, user.Username)
	}()

	storage := New(db)
	id, err := storage.Create(context.Background(), makeEvent(user.ID))
	require.NoError(t, err)
	require.NotZero(t, id)

	deleteEvent(t, db, id)
}

func TestEventRepository_Get(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	user := createUser(t, db)
	defer func() {
		deleteUserByUserName(t, db, user.Username)
	}()

	expEvent := makeEvent(user.ID)
	id := createEvent(t, db, expEvent)
	defer func() {
		deleteEvent(t, db, id)
	}()

	expEvent.ID = id
	storage := New(db)

	tests := []struct {
		name string
		id   int64
		exp  entity.Event
		err  error
	}{
		{
			name: "event not found",
			id:   id + 1,
			err:  store.ErrEventNotFound,
		},
		{
			name: "success",
			id:   id,
			exp:  expEvent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEvent, err := storage.Get(context.Background(), tt.id)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.exp, gotEvent)
		})
	}
}

func TestEventRepository_Update(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	user := createUser(t, db)
	defer func() {
		deleteUserByUserName(t, db, user.Username)
	}()

	event := makeEvent(user.ID)
	id := createEvent(t, db, event)
	defer func() {
		deleteEvent(t, db, id)
	}()

	event.Title = "1"
	event.Description = "2"
	event.Duration = 19 * time.Second
	event.BeforeStartNotice = 20 * time.Second
	event.Time = dateTime.Add(1 * time.Hour)

	storage := New(db)

	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "event not found",
			id:   id + 1,
			err:  store.ErrEventNotFound,
		},
		{
			name: "success",
			id:   id,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.id)
			event.ID = tt.id
			err := storage.Update(context.Background(), event)
			require.ErrorIs(t, err, tt.err)
		})
	}

}

func TestEventRepository_Delete(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	user := createUser(t, db)
	defer func() {
		deleteUserByUserName(t, db, user.Username)
	}()

	event := makeEvent(user.ID)
	id := createEvent(t, db, event)

	storage := New(db)

	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "event not found",
			id:   id + 1,
			err:  store.ErrEventNotFound,
		},
		{
			name: "success",
			id:   id,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Delete(context.Background(), tt.id)
			require.ErrorIs(t, err, tt.err)
		})
	}

}
