package storage

type Config struct {
	Driver   string `env:"DRIVER" yaml:"driver" envDefault:"memory"`
	User     string `env:"USER" yaml:"user"`
	Password string `env:"PASSWORD" yaml:"password"`
	DB       string `env:"DB" yaml:"db"`
	Host     string `env:"HOST" envDefault:"localhost" yaml:"host"`
	Port     int    `env:"PORT" envDefault:"3306" yaml:"port"`
}
