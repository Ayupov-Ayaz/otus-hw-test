package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type publisher struct {
	ch *amqp.Channel
}

func newPublisher(ch *amqp.Channel) *publisher {
	return &publisher{ch: ch}
}

func (p *publisher) Publish(event string, body []byte) error {
	q, err := p.ch.QueueDeclare(event, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = p.ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
