package internalhttp

import (
	"time"

	goFiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LogRequestMiddleware(logger *zap.Logger) goFiber.Handler {
	return func(ctx *goFiber.Ctx) error {
		start := time.Now()

		defer func() {
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
		}()

		return ctx.Next()
	}
}
