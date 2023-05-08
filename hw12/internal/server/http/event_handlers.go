package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	goFiber "github.com/gofiber/fiber/v2"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"
)

const eventID = "id"

var (
	ErrEventIDNotFound = errors.New("event id not found")
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, e entity.Event) (int64, error)
	UpdateEvent(ctx context.Context, e entity.Event) error
	DeleteEvent(ctx context.Context, id int64) error
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
