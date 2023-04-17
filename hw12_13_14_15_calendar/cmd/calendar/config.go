package main

import "strconv"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	HTTP   HTTPServerConf
	// TODO
}

func DefaultConfig() Config {
	return Config{
		Logger: DefaultLoggerConf(),
		HTTP:   DefaultHTTPServerConf(),
	}
}

type LoggerConf struct {
	Level string
}

func DefaultLoggerConf() LoggerConf {
	return LoggerConf{
		Level: "debug",
	}
}

type HTTPServerConf struct {
	Port int
}

func DefaultHTTPServerConf() HTTPServerConf {
	return HTTPServerConf{
		Port: 8080,
	}
}

func (c HTTPServerConf) PortToString() string {
	return strconv.Itoa(c.Port)
}

func NewConfig() Config {
	cfg := DefaultConfig()
	// todo: unmarshal environment variables
	return cfg
}
