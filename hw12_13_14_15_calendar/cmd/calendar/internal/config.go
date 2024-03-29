package internal

import (
	"errors"
	"fmt"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal/configs/grpc"
	"os"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal/configs/logger"

	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal/configs/http"
	"github.com/ayupov-ayaz/otus-wh-test/hw12_13_14_15_calendar/cmd/calendar/internal/configs/storage"

	yaml3 "gopkg.in/yaml.v3"

	env8 "github.com/caarlos0/env/v8"
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

func unmarshalYamlFile(yamlFile string) func(cfg *Config) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return func(cfg *Config) error {
			if errors.Is(err, os.ErrNotExist) {
				// Если файла нет, то просто возвращаем пустую функцию,
				return nil
			}

			// Если файл есть, но не удалось его прочитать, возвращаем ошибку.
			return err
		}
	}

	return func(cfg *Config) error {
		return unmarshalYaml(data)(cfg)
	}
}

func unmarshalEnv(cfg *Config) error {
	opts := env8.Options{
		Prefix: envPrefix,
	}

	return env8.ParseWithOptions(cfg, opts)
}

// NewConfig returns a new Config.
// order: ENV, YAML
func NewConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	if err := unmarshalEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal env: %w", err)
	}

	if err := unmarshalYamlFile(configFile)(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
