package internal

import (
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/grpc"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/logger"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/http"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/configs/settings/storage"

	yaml3 "gopkg.in/yaml.v3"
)

const (
	envPrefix = "CALENDAR_"
)

type Config struct {
	Logger  logger.Config  `envPrefix:"LOGGER_"`
	HTTP    http.Config    `envPrefix:"HTTP_"`
	GRPC    grpc.Config    `envPrefix:"GRPC_"`
	Storage storage.Config `envPrefix:"STORAGE_"`
}

func unmarshalYaml(data []byte) func(cfg *Config) error {
	return func(cfg *Config) error {
		if err := yaml3.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("failed to unmarshal yaml: %w", err)
		}

		return nil
	}
}

// NewConfig returns a new Config.
// order: ENV, YAML
func NewConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	if err := settings.UnmarshalEnv(envPrefix, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal env: %w", err)
	}

	if err := settings.UnmarshalYamlFile(configFile, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
