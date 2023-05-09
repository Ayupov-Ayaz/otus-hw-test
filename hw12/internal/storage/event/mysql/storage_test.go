//go:build integration
// +build integration

package mysql

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	store "github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	testDSN  string
	dateTime entity.MyTime
)

func TestMain(m *testing.M) {
	testDSN = test.GetMysqlTestDSN()
	var err error
	dt, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	if err != nil {
		panic(err)
	}

	dateTime = entity.MyTime(dt)

	os.Exit(m.Run())
}

func getConnection(t *testing.T) *sqlx.DB {
	return test.MysqlConnection(t, testDSN)
}

func deleteEvent(t *testing.T, db *sqlx.DB, id int64) {
	res, err := db.Exec(deleteQuery, id)
	require.NoError(t, err)
	rows, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)
}

func createEvent(t *testing.T, db *sqlx.DB, e entity.Event) int64 {
	res, err := db.Exec(createQuery, e.Title, e.UserID, e.Description, e.DateTime.Time(),
		e.DurationInSeconds())
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return id
}

func makeEvent(userID int64) entity.Event {
	duration := entity.NewSecondsDuration(5)
	title := strconv.Itoa(time.Now().Nanosecond())
	return entity.NewEvent(title, "desc", userID, dateTime, duration, nil)
}

func TestEventRepository_Create(t *testing.T) {
	makeEvent(1)
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	storage := New(db)
	id, err := storage.Create(context.Background(), makeEvent(1))
	require.NoError(t, err)
	require.NotZero(t, id)

	deleteEvent(t, db, id)
}

func TestEventRepository_Get(t *testing.T) {
	db := getConnection(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	expEvent := makeEvent(1)
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

	event := makeEvent(1)
	id := createEvent(t, db, event)
	defer func() {
		deleteEvent(t, db, id)
	}()

	event.Title = "1"
	event.Description = "2"
	event.Duration = entity.NewSecondsDuration(19)
	event.DateTime = entity.MyTime(dateTime.Time().Add(1 * time.Hour))

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

	event := makeEvent(1)
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
