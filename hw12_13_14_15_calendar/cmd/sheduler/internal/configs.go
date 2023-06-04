package internal

import (
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/event"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/intervals"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/logger"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/queue"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/storage"
)

const envPrefix = "QUEUE_"

type Config struct {
	Logger    logger.Config    `envPrefix:"LOGGER_"`
	Event     event.Config     `envPrefix:"EVENT_"`
	Queue     queue.Config     `envPrefix:"QUEUE_"`
	Storage   storage.Config   `envPrefix:"STORAGE_"`
	Intervals intervals.Config `envPrefix:"INTERVALS_"`
}

func NewConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	if err := settings.UnmarshalEnv(envPrefix, cfg); err != nil {
		return nil, err
	}

	if err := settings.UnmarshalYamlFile(configFile, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
