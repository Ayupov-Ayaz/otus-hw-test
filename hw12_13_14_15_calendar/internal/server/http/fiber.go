package internalhttp

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func errorHandler(logger *zap.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		data, marshalErr := jsoniter.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})

		if marshalErr != nil {
			logger.Error("marshaling failed", zap.Error(marshalErr))
		}

		if err := ctx.Status(fiber.StatusInternalServerError).Send(data); err != nil {
			logger.Error("sending response failed", zap.Error(err))
		}

		return nil
	}
}

func newFiber(logger *zap.Logger, version string) *fiber.App {
	f := fiber.New(fiber.Config{
		AppName:      "Calendar " + version,
		ErrorHandler: errorHandler(logger),
	})

	f.Use(LogRequestMiddleware(logger))

	return f
}
