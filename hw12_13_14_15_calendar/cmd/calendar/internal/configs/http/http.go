package http

import "strconv"

type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"8080"`
}

func (c Config) PortToString() string {
	return strconv.Itoa(c.Port)
}

func (c Config) Addr() string {
	return c.Host + ":" + c.PortToString()
}
