package event

type Config struct {
	Name string `env:"NAME" envDefault:"events"`
}
