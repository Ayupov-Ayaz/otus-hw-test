package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/caarlos0/env/v8"
)

const (
	envPrefix = "CALENDAR_"
)

type Config struct {
	Logger  LoggerConf     `envPrefix:"LOGGER_"`
	HTTP    HTTPServerConf `envPrefix:"HTTP_"`
	Storage StorageConf    `envPrefix:"STORAGE_"`
}

type LoggerConf struct {
	Level string ` env:"LEVEL" envDefault:"debug"`
}

type HTTPServerConf struct {
	Port int ` env:"PORT" envDefault:"8080"`
}

func (c HTTPServerConf) PortToString() string {
	return strconv.Itoa(c.Port)
}

func unmarshalYaml(data []byte) func(cfg *Config) error {
	return func(cfg *Config) error {
		if err := yaml.Unmarshal(data, cfg); err != nil {
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
	opts := env.Options{
		Prefix: envPrefix,
	}

	return env.ParseWithOptions(cfg, opts)
}

// NewConfig returns a new Config.
// 1. ENV
// 2. YAML
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := unmarshalEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal env: %w", err)
	}

	if err := unmarshalYamlFile(configFile)(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return cfg, nil
}
