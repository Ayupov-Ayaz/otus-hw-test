package api

import (
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func castGRPCEvent(event *Event) entity.Event {
	e := entity.Event{
		Title:        event.Title,
		Description:  event.Description,
		UserID:       event.UserId,
		DateTime:     entity.MyTime(event.Time.AsTime()),
		Duration:     entity.Duration(event.Duration.AsDuration()),
		Notification: entity.Duration(event.Notification.AsDuration()),
	}

	if event.Id != nil {
		v := *event.Id
		e.ID = v
	}

	return e
}

func castEvent(event entity.Event) *Event {
	v := event.ID

	return &Event{
		Id:           &v,
		Title:        event.Title,
		Description:  event.Description,
		UserId:       event.UserID,
		Time:         timestamppb.New(event.DateTime.Time()),
		Duration:     durationpb.New(time.Duration(event.Duration)),
		Notification: durationpb.New(time.Duration(event.Notification)),
	}
}

func castEvents(events []entity.Event) []*Event {
	result := make([]*Event, len(events))
	for i, event := range events {
		result[i] = castEvent(event)
	}

	return result
}
