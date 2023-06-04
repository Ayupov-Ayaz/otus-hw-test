package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sheduler/internal"
	app "github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/app/sheduler"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/signals"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/internal/storage/connect"
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

	cfg, err := internal.NewConfig(configFile)
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

	storage, err := connect.NewStorage(cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	publisher := queue.Publisher()
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	publisherConfigs := []config{
		{
			callInterval: cfg.Intervals.Day,
			getDuration:  24 * time.Hour,
			name:         "day",
		},
		{
			callInterval: cfg.Intervals.Week,
			getDuration:  7 * 24 * time.Hour,
			name:         "week",
		},
		{
			callInterval: cfg.Intervals.Month,
			getDuration:  30 * 24 * time.Hour,
			name:         "month",
		},
	}

	pub := app.New(storage, publisher, logg)
	for _, pCfg := range publisherConfigs {
		callback := pub.Sender(cfg.Event.Name, pCfg.getDuration)

		go pub.Start(notifyCtx.Done(), pCfg.name, callback, pCfg.callInterval)
	}

	pub.Remover(notifyCtx.Done(), cfg.Intervals.Remove)

	<-notifyCtx.Done()

	return nil
}

type config struct {
	callInterval time.Duration
	getDuration  time.Duration
	name         string
}
