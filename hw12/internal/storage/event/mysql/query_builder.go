package mysql

import (
	"strconv"
	"strings"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

const mysqlLayout = "2006-01-02 15:04:05"

func mySQLTimeFormat(time time.Time) string {
	return time.Format(mysqlLayout)
}

func parseMySQLTime(str string) (time.Time, error) {
	return time.Parse(mysqlLayout, str)
}

type QueryBuilder struct{}

func (QueryBuilder) createNotificationQuery(eventID int64, notification []entity.Duration) string {
	var b strings.Builder
	strEventID := strconv.FormatInt(eventID, 10)
	b.WriteString("INSERT INTO notifications (event_id, before_start_notice_sec) VALUES ")
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
	b.WriteString(";")
	return b.String()
}

func (QueryBuilder) updateEventQuery(event entity.Event) (string, bool) {
	var (
		b    strings.Builder
		open bool
	)

	b.WriteString("UPDATE events SET ")
	add := func(field, value, wrap string) {
		if open {
			b.WriteString(", ")
		}
		open = true
		b.WriteString(field)
		b.WriteString(" = ")
		if wrap != "" {
			b.WriteString(wrap)
		}
		b.WriteString(value)
		if wrap != "" {
			b.WriteString(wrap)
		}
	}

	if event.Title != "" {
		add("title", event.Title, "'")
	}
	if event.Description != "" {
		add("description", event.Description, "'")
	}
	if !event.DateTime.IsEmpty() {
		add("time", mySQLTimeFormat(event.DateTime.Time()), "'")
	}
	if !event.Duration.IsEmpty() {
		add("duration_sec", strconv.Itoa(event.Duration.DurationInSec()), "")
	}
	if event.UserID != 0 {
		add("user_id", strconv.FormatInt(event.UserID, 10), "")
	}

	b.WriteString(" WHERE id = ")
	b.WriteString(strconv.FormatInt(event.ID, 10))
	b.WriteString(";")
	return b.String(), open
}
