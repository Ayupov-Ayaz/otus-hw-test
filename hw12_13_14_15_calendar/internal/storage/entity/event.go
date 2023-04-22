package entity

import "time"

type Event struct {
	ID                int
	Description       string
	Title             string        `required:"true"`
	UserID            string        `required:"true"`
	Time              time.Time     `required:"true"`
	Duration          time.Duration `required:"true"`
	BeforeStartNotice time.Duration
}

func NewEvent(description string, title string, userID string,
	time time.Time, duration time.Duration, beforeStartNotice time.Duration) Event {
	return Event{
		Description:       description,
		Title:             title,
		UserID:            userID,
		Time:              time,
		Duration:          duration,
		BeforeStartNotice: beforeStartNotice,
	}
}
