package entity

import "time"

type Event struct {
	ID                int64
	Description       string
	Title             string        `required:"true"`
	UserID            int64         `required:"true"`
	Time              time.Time     `required:"true"`
	Duration          time.Duration `required:"true"`
	BeforeStartNotice time.Duration
}

func NewEvent(title string, description string, userID int64,
	time time.Time, duration time.Duration, beforeStartNotice time.Duration,
) Event {
	return Event{
		Description:       description,
		Title:             title,
		UserID:            userID,
		Time:              time,
		Duration:          duration,
		BeforeStartNotice: beforeStartNotice,
	}
}

func (e Event) DurationInSeconds() int {
	return int(e.Duration.Seconds())
}

func (e Event) BeforeStartNoticeInSeconds() int {
	return int(e.BeforeStartNotice.Seconds())
}
