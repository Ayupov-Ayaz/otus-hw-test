package internalhttp

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	app    Application
	logger *zap.Logger
	f      *fiber.App
}

type Application interface { // TODO
}

func NewServer(logger *zap.Logger, app Application, version string) *Server {
	fiberApp := newFiber(logger, version)

	srv := &Server{
		app:    app,
		logger: logger,
		f:      fiberApp,
	}

	fiberApp.Get("/hello", srv.HelloWorld)

	return srv
}

func (s *Server) Start(_ context.Context, address string) error {
	return s.f.Listen(address)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping server")
	return s.f.ShutdownWithContext(ctx)
}

func (s *Server) HelloWorld(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World 👋!")
}
