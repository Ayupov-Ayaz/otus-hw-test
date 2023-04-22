package main

import "time"

const (
	MySQL  = "mysql"
	Memory = "memory"
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
