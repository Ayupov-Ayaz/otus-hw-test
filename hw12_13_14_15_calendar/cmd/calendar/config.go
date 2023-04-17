package main

import (
	"strconv"

	"github.com/caarlos0/env/v8"
)

const envPrefix = "CALENDAR_"

type Config struct {
	Logger LoggerConf     `envPrefix:"LOGGER_"`
	HTTP   HTTPServerConf `envPrefix:"HTTP_"`
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

func NewConfig() (*Config, error) {
	cfg := &Config{}

	opts := env.Options{
		Prefix: envPrefix,
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	return cfg, nil
}
