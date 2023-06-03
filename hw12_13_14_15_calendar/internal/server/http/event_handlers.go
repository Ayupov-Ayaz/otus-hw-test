package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	goFiber "github.com/gofiber/fiber/v2"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/entity"
)

const (
	eventID   = "id"
	userID    = "user_id"
	startTime = "start"
)

var (
	ErrEventIDNotFound = errors.New("event id not found")
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, e entity.Event) (int64, error)
	UpdateEvent(ctx context.Context, e entity.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsByWeek(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
	GetEventsByMonth(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
	GetEventsByDay(ctx context.Context, userID int64, date time.Time) ([]entity.Event, error)
}

type EventHandlers struct {
	logger *zap.Logger
	app    EventUseCase
}

func NewEventHandlers(logger *zap.Logger, app EventUseCase) *EventHandlers {
	return &EventHandlers{
		logger: logger,
		app:    app,
	}
}

func (e *EventHandlers) Register(router goFiber.Router) {
	router.Post("/", e.CreateEvent)
	router.Post("/:"+eventID, e.UpdateEvent)
	router.Delete("/:"+eventID, e.DeleteEvent)
	router.Get("/user/:"+userID+"/week/:"+startTime, e.GetEventsByWeek)
	router.Get("/user/:"+userID+"/day/:"+startTime, e.GetEventsByDay)
	router.Get("/user/:"+userID+"/month/:"+startTime, e.GetEventsByMonth)
}
func unmarshalEvent(logger *zap.Logger, data []byte) (entity.Event, error) {
	var event entity.Event
	if err := jsoniter.Unmarshal(data, &event); err != nil {
		logger.Error("failed to unmarshal event",
			zap.ByteString("body", data),
			zap.Error(err))
		return entity.Event{}, err
	}

	return event, nil
}

func createEventResponse(id int64) []byte {
	return []byte(`{"id":` + strconv.FormatInt(id, 10) + `}`)
}

func getEventID(ctx *goFiber.Ctx, logger *zap.Logger) (int64, error) {
	params := ctx.AllParams()
	id, ok := params[eventID]
	if !ok {
		err := ErrEventIDNotFound
		logger.Error("failed to get event id", zap.Error(err))
		return 0, err
	}

	resp, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to parse event id", zap.Error(err))
		return 0, err
	}

	return resp, nil
}

func (e *EventHandlers) CreateEvent(ctx *goFiber.Ctx) error {
	event, err := unmarshalEvent(e.logger, ctx.Body())
	if err != nil {
		return err
	}

	id, err := e.app.CreateEvent(ctx.Context(), event)
	if err != nil {
		return err
	}

	sendJson(ctx, goFiber.StatusCreated, createEventResponse(id))

	return nil
}

func (e *EventHandlers) UpdateEvent(ctx *goFiber.Ctx) error {
	event, err := unmarshalEvent(e.logger, ctx.Body())
	if err != nil {
		return err
	}

	id, err := getEventID(ctx, e.logger)
	if err != nil {
		return err
	}

	event.ID = id

	if err := e.app.UpdateEvent(ctx.Context(), event); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}

func (e *EventHandlers) DeleteEvent(ctx *goFiber.Ctx) error {
	id, err := getEventID(ctx, e.logger)
	if err != nil {
		return err
	}

	if err := e.app.DeleteEvent(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}

func getUserID(ctx *goFiber.Ctx) (int64, error) {
	userID, ok := ctx.AllParams()[userID]
	if !ok {
		err := errors.New("failed to get userID from context")
		return 0, err
	}

	resp, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse userID: %w", err)
	}

	return resp, nil
}

func getDate(ctx *goFiber.Ctx) (time.Time, error) {
	date, ok := ctx.AllParams()[startTime]
	if !ok {
		err := errors.New("failed to get date from context")
		return time.Time{}, err
	}

	resp, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return resp, nil
}

func (e *EventHandlers) getDateAndUserID(ctx *goFiber.Ctx) (date time.Time, userID int64, err error) {
	userID, err = getUserID(ctx)
	if err != nil {
		e.logger.Error("failed to get userID", zap.Error(err))
		return time.Time{}, 0, err
	}

	date, err = getDate(ctx)
	if err != nil {
		e.logger.Error("failed to get date", zap.Error(err))
		return time.Time{}, 0, err
	}

	return date, userID, nil
}

func (e *EventHandlers) sendEvents(ctx *goFiber.Ctx, events []entity.Event) error {
	resp, err := jsoniter.Marshal(struct {
		Events []entity.Event `json:"events"`
	}{
		Events: events,
	})

	if err != nil {
		e.logger.Error("failed to marshal events", zap.Error(err))
		return err
	}

	sendJson(ctx, goFiber.StatusOK, resp)

	return nil
}

func (e *EventHandlers) GetEventsByDay(ctx *goFiber.Ctx) error {
	date, userID, err := e.getDateAndUserID(ctx)
	if err != nil {
		return err
	}

	events, err := e.app.GetEventsByDay(ctx.Context(), userID, date)
	if err != nil {
		return err
	}

	return e.sendEvents(ctx, events)
}

func (e *EventHandlers) GetEventsByWeek(ctx *goFiber.Ctx) error {
	date, userID, err := e.getDateAndUserID(ctx)
	if err != nil {
		return err
	}

	events, err := e.app.GetEventsByWeek(ctx.Context(), userID, date)
	if err != nil {
		return err
	}

	return e.sendEvents(ctx, events)
}

func (e *EventHandlers) GetEventsByMonth(ctx *goFiber.Ctx) error {
	date, userID, err := e.getDateAndUserID(ctx)
	if err != nil {
		return err
	}

	events, err := e.app.GetEventsByMonth(ctx.Context(), userID, date)
	if err != nil {
		return err
	}

	return e.sendEvents(ctx, events)
}
