package internalhttp

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/ayupov-ayaz/otus-wh-test/hw12/internal/storage/entity"

	goFiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Application interface {
	CreateEvent(ctx context.Context, e entity.Event) error
}

type Server struct {
	app    Application
	logger *zap.Logger
	f      *goFiber.App
}

func NewServer(logger *zap.Logger, app Application, version string) *Server {
	fiberApp := newFiber(logger, version)

	srv := &Server{
		app:    app,
		logger: logger,
		f:      fiberApp,
	}

	group := fiberApp.Group("/event")
	group.Post("/", srv.CreateEvent)

	return srv
}

func (s *Server) Start(_ context.Context, address string) error {
	return s.f.Listen(address)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping server")
	return s.f.ShutdownWithContext(ctx)
}

func (s *Server) CreateEvent(ctx *goFiber.Ctx) error {
	var event entity.Event
	if err := jsoniter.Unmarshal(ctx.Body(), &event); err != nil {
		s.logger.Error("failed to unmarshal event",
			zap.ByteString("body", ctx.Body()),
			zap.Error(err))
		return err
	}

	if err := s.app.CreateEvent(ctx.Context(), event); err != nil {
		return err
	}

	return ctx.SendStatus(goFiber.StatusCreated)
}
