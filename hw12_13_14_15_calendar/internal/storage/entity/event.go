package entity

import "time"

type Event struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title" validate:"required"`
	Description  string   `json:"description"`
	UserID       int64    `json:"user_id" validate:"required"`
	DateTime     MyTime   `json:"time" validate:"required"`
	Duration     Duration `json:"duration" validate:"required"`
	Notification Duration `json:"notification" validate:"required"`
}

func NewEvent(
	title, description string,
	userID int64,
	time MyTime,
	duration, notifications Duration,
) Event {
	return Event{
		Description:  description,
		Title:        title,
		UserID:       userID,
		DateTime:     time,
		Duration:     duration,
		Notification: notifications,
	}
}

func (e Event) DurationInSeconds() int {
	return e.Duration.DurationInSec()
}

func (e Event) EventDate() time.Time {
	return e.DateTime.Time()
}

func (e *Event) Reset() {
	*e = Event{}
}

func (e Event) NotificationInSec() int {
	return e.Notification.DurationInSec()
}
