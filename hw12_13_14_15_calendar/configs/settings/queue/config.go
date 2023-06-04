package queue

type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"5672"`
	User string `env:"USER"`
	Pass string `env:"PASS"`
}
