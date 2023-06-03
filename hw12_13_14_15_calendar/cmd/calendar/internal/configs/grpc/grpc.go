package grpc

import (
	"net"
	"strconv"
)

type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"9000"`
}

func (c Config) PortToString() string {
	return strconv.Itoa(c.Port)
}

func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, c.PortToString())
}
