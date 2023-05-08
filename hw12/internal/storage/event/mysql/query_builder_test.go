package mysql

import (
	"testing"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder(t *testing.T) {
	const exp = "INSERT INTO notification (event_id, before_start_notice) VALUES (1, 1), (1, 3), (1, 5)"

	notifications := []entity.Duration{
		entity.NewSecondsDuration(1),
		entity.NewSecondsDuration(3),
		entity.NewSecondsDuration(5),
	}

	got := createNotificationQuery(1, notifications)
	require.Equal(t, exp, got)
}
