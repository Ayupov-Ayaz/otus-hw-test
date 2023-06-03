package main

import (
	"flag"
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sender/app"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sender/internal/config"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/signals"
	"go.uber.org/zap"
)

const defaultConfigFile = "/etc/calendar/config.toml"

var configFile string

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "Path to configuration file")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	logg := logger.New(cfg.Logger.Level)
	defer func() {
		_ = logg.Sync()
	}()

	logg.Info("using config file", zap.String("path", configFile))

	logg.Info("using rabbit queue",
		zap.String("host", cfg.Queue.Host),
		zap.Int("port", cfg.Queue.Port),
		zap.String("event_name", cfg.Event.Name))

	notifyCtx, cancel := signals.NotifyCtx()
	defer cancel()

	queue, err := rabbit.New(rabbit.Config(cfg.Queue))
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	events, err := queue.Consumer().Consume(cfg.Event.Name)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	logg.Info("start consuming")

	app.Start(notifyCtx.Done(), events, logg)

	return nil
}
