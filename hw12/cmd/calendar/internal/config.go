package internal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	yaml3 "gopkg.in/yaml.v3"

	env8 "github.com/caarlos0/env/v8"
)

const (
	envPrefix = "CALENDAR_"
)

type Timeouts struct {
	Read time.Duration `env:"READ_TIMEOUT" envDefault:"5s" yaml:"read"`
}

type StorageConf struct {
	Driver   string   `env:"DRIVER" yaml:"driver" envDefault:"memory"`
	User     string   `env:"USER" yaml:"user"`
	Password string   `env:"PASSWORD" yaml:"password"`
	DB       string   `env:"Storage" yaml:"db"`
	Host     string   `env:"HOST" envDefault:"localhost" yaml:"host"`
	Port     int      `env:"PORT" envDefault:"3306" yaml:"port"`
	Timeouts Timeouts `envPrefix:"TIMEOUTS_" yaml:"timeouts"`
}

func (s StorageConf) IsMemoryStorage() bool {
	return s.Driver == Memory
}

type Config struct {
	Logger  LoggerConf     `envPrefix:"LOGGER_"`
	HTTP    HTTPServerConf `envPrefix:"HTTP_"`
	Storage StorageConf    `envPrefix:"STORAGE_"`
}

type LoggerConf struct {
	Level string `env:"LEVEL" envDefault:"debug"`
}

type HTTPServerConf struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"8080"`
}

func (c HTTPServerConf) PortToString() string {
	return strconv.Itoa(c.Port)
}

func (c HTTPServerConf) Addr() string {
	return c.Host + ":" + c.PortToString()
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
