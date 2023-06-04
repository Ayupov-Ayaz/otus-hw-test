package intervals

import "time"

type Config struct {
	Day    time.Duration `env:"DAY" envDefault:"1m"`
	Week   time.Duration `env:"WEEK" envDefault:"1h"`
	Month  time.Duration `env:"MONTH" envDefault:"24h"`
	Remove time.Duration `env:"REMOVE" envDefault:"24h"`
}
