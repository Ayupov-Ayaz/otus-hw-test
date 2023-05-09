package mysql

import (
	"testing"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_createNotification(t *testing.T) {
	qb := QueryBuilder{}
	const exp = "INSERT INTO notifications (event_id, before_start_notice_sec) VALUES (1, 1), (1, 3), (1, 5);"

	notifications := []entity.Duration{
		entity.NewSecondsDuration(1),
		entity.NewSecondsDuration(3),
		entity.NewSecondsDuration(5),
	}

	got := qb.createNotificationQuery(1, notifications)
	require.Equal(t, exp, got)
}

func ParseTime(dateTime string) (entity.MyTime, error) {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return entity.MyTime{}, err
	}

	return entity.NewMyTime(t), nil
}

func TestQueryBuilder_updateEvent(t *testing.T) {
	const (
		dateTime      = "2100-04-05T12:01:01Z"
		mySqlDateTime = "2100-04-05 12:01:01"
	)

	dt, err := ParseTime(dateTime)
	require.NoError(t, err)

	dur := entity.NewSecondsDuration(23)
	qb := QueryBuilder{}
	tests := []struct {
		exp   string
		event entity.Event
		ok    bool
	}{
		{
			exp: "UPDATE events SET title = 'test', description = 'test', time = '" +
				mySqlDateTime + "', duration_sec = 23, before_notice_sec = 23, user_id = 45 WHERE id = 1;",
			ok: true,
			event: entity.Event{
				ID:           1,
				Title:        "test",
				Description:  "test",
				DateTime:     dt,
				Duration:     dur,
				UserID:       45,
				Notification: entity.NewSecondsDuration(23),
			},
		},
		{
			exp: "UPDATE events SET  WHERE id = 0;",
			ok:  false,
		},
		{
			exp: "UPDATE events SET title = 'qwerty' WHERE id = 2;",
			ok:  true,
			event: entity.Event{
				Title: "qwerty",
				ID:    2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.exp, func(t *testing.T) {
			got, ok := qb.updateEventQuery(tt.event)
			require.Equal(t, tt.exp, got)
			require.Equal(t, tt.ok, ok)
		})
	}
}
