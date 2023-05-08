package mysql

import (
	"strconv"
	"strings"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

const insertNotificationQuery = "INSERT INTO notifications (event_id, before_start_notice_sec) VALUES "

func createNotificationQuery(eventID int64, notification []entity.Duration) string {
	var b strings.Builder
	strEventID := strconv.FormatInt(eventID, 10)
	b.WriteString(insertNotificationQuery)
	count := len(notification)
	for i, n := range notification {
		b.WriteString("(")
		b.WriteString(strEventID)
		b.WriteString(", ")
		b.WriteString(strconv.Itoa(n.DurationInSec()))
		b.WriteString(")")
		if i != count-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}
