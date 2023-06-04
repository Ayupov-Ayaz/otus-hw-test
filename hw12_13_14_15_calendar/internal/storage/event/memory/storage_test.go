package memory

import (
	"context"
	"testing"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func makeEvent() entity.Event {
	return entity.NewEvent("title", "desc",
		12, entity.MyTime(time.Now()), entity.Duration(2*time.Second),
		entity.NewSecondsDuration(5))
}

func TestStorage_Create(t *testing.T) {
	storage := New()
	ctx := context.Background()
	e := makeEvent()
	for i := 0; i < 10; i++ {
		expID := int64(i + 1)
		id, err := storage.Create(ctx, e)
		require.NoError(t, err)
		require.Equal(t, expID, id)
	}
}

func TestStorage_Update(t *testing.T) {
	const id = 12
	e := makeEvent()
	e.ID = id
	store := New()
	store.events[id] = e

	ctx := context.Background()
	e.Title = "title 2"
	e.Description = "desc 2"

	tests := []struct {
		name string
		err  error
		e    entity.Event
	}{
		{
			name: "event not found",
		},
		{
			name: "success",
			e:    e,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Update(ctx, tt.e)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.e, store.events[tt.e.ID])
		})
	}
}

func TestStorage_Get(t *testing.T) {
	const (
		id     = 11
		userID = 134
	)
	e := makeEvent()
	e.ID = id
	e.UserID = userID

	e.DateTime = entity.MyTime(time.Now().Add(24 * time.Hour))
	store := New()
	store.events[id] = e
	ctx := context.Background()

	tests := []struct {
		name    string
		exp     []entity.Event
		addTime time.Duration
		err     error
	}{
		{
			name:    "event not found",
			addTime: 23 * time.Hour,
		},
		{
			name:    "success",
			addTime: 24 * time.Hour,
			exp:     []entity.Event{e},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetEventsForDates(ctx, time.Now(), time.Now().Add(tt.addTime))
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	const id = 33
	e := makeEvent()
	e.ID = id
	store := New()
	store.events[id] = e
	ctx := context.Background()

	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "event not found",
			id:   12,
		},
		{
			name: "success",
			id:   id,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(ctx, tt.id)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
