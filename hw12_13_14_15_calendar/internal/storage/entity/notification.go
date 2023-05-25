package entity

import "time"

type Notification struct {
	EventID           int64
	BeforeStartNotice time.Duration
}

func NewNotification(eventID int64, beforeStartNotice time.Duration) Notification {
	return Notification{
		EventID:           eventID,
		BeforeStartNotice: beforeStartNotice,
	}
}

func (n Notification) BeforeStartNoticeInSeconds() int {
	return int(n.BeforeStartNotice.Seconds())
}
