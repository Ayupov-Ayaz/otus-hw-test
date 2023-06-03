package config

import (
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sender/internal/config/event"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sender/internal/config/logger"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/sender/internal/config/queue"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/parser"
)

const envPrefix = "QUEUE_"

type Config struct {
	Logger logger.Config `envPrefix:"LOGGER_"`
	Event  event.Config  `envPrefix:"EVENT_"`
	Queue  queue.Config  `envPrefix:"QUEUE_"`
}

func NewConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	if err := parser.UnmarshalEnv(envPrefix, cfg); err != nil {
		return nil, err
	}

	if err := parser.UnmarshalYamlFile(configFile, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
