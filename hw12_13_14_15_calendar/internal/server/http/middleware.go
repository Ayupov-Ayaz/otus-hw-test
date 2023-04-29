package internalhttp

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func LogRequestMiddleware(logger *zap.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		if err := ctx.Next(); err != nil {
			return err
		}

		measure := time.Since(start)

		logger.Info("request",
			zap.String("protocol", ctx.Protocol()),
			zap.String("ip", ctx.IP()),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.OriginalURL()),
			zap.Int("status", ctx.Response().StatusCode()),
			zap.Duration("latency", measure),
			zap.String("user-agent", ctx.Get("User-Agent")),
		)

		return nil
	}
}
