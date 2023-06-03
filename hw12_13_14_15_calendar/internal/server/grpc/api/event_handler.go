package api

import (
	"context"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
	"go.uber.org/zap"
	"time"
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, e entity.Event) (int64, error)
	UpdateEvent(ctx context.Context, e entity.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsByWeek(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
	GetEventsByMonth(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
	GetEventsByDay(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
}

type EventHandler struct {
	app  EventUseCase
	logg *zap.Logger
	UnimplementedEventServiceServer
}

func NewEventHandler(app EventUseCase, logg *zap.Logger) *EventHandler {
	return &EventHandler{
		app:  app,
		logg: logg,
	}
}

func (e *EventHandler) CreateEvent(ctx context.Context, request *CreateEventRequest) (*CreateEventResponse, error) {
	event := castGRPCEvent(request.Event)
	id, err := e.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	resp := &CreateEventResponse{
		Id: id,
	}

	return resp, nil
}

func (e *EventHandler) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*UpdateEventResponse, error) {
	event := castGRPCEvent(request.Event)
	if err := e.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}

	return &UpdateEventResponse{}, nil
}

func (e *EventHandler) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*DeleteEventResponse, error) {
	if err := e.app.DeleteEvent(ctx, request.Id); err != nil {
		return nil, err
	}

	return &DeleteEventResponse{}, nil
}

func (e *EventHandler) getEvents(ctx context.Context, userID int64, date time.Time, getFunc func(context.Context, int64, time.Time) ([]entity.Event, error)) (*GetEventsResponse, error) {
	events, err := getFunc(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	resp := &GetEventsResponse{
		Events: castEvents(events),
	}

	return resp, nil
}

func (e *EventHandler) GetEventByDay(ctx context.Context, request *GetEventsRequest) (*GetEventsResponse, error) {
	return e.getEvents(ctx, request.UserId, request.Time.AsTime(), e.app.GetEventsByDay)
}

func (e *EventHandler) GetEventByWeek(ctx context.Context, request *GetEventsRequest) (*GetEventsResponse, error) {
	return e.getEvents(ctx, request.UserId, request.Time.AsTime(), e.app.GetEventsByWeek)
}

func (e *EventHandler) GetEventByMonth(ctx context.Context, request *GetEventsRequest) (*GetEventsResponse, error) {
	return e.getEvents(ctx, request.UserId, request.Time.AsTime(), e.app.GetEventsByMonth)
}
