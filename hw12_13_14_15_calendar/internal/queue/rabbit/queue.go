package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(event string, body []byte) error
}

type Consumer interface {
	Consume(event string) (<-chan []byte, error)
}

type Config struct {
	Host string
	Port int
	User string
	Pass string
}

func (c Config) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.User, c.Pass, c.Host, c.Port)
}

type Queue struct {
	event string
	ch    *amqp.Channel
}

func New(config Config) (*Queue, error) {
	conn, err := amqp.Dial(config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	q := &Queue{ch: ch}

	return q, nil
}

func (q *Queue) Publisher() Publisher {
	return newPublisher(q.ch)
}

func (q *Queue) Consumer() Consumer {
	return newConsumer(q.ch)
}
