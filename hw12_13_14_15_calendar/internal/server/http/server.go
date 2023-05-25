package internalhttp

import (
	"context"

	goFiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	f      *goFiber.App
}

func NewServer(logger *zap.Logger, version string) *Server {
	fiberApp := newFiber(logger, version)

	srv := &Server{
		logger: logger,
		f:      fiberApp,
	}

	return srv
}

func sendJson(ctx *goFiber.Ctx, status int, raw []byte) {
	ctx.Status(status)
	ctx.Response().Header.SetContentType(goFiber.MIMEApplicationJSON)
	ctx.Response().SetBodyRaw(raw)
}

func (s *Server) Start(_ context.Context, address string) error {
	return s.f.Listen(address)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping server")
	return s.f.ShutdownWithContext(ctx)
}

func (s *Server) Register(handlers *EventHandlers) {
	handlers.Register(s.f.Group("/event"))
}
