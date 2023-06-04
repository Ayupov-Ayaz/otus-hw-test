package rabbit

import (
	"fmt"
	"runtime"

	"github.com/streadway/amqp"
)

const (
	autoAck = true
	// single consumer
	exclusive = true
	noLocal   = false
	noWait    = false
)

type consumer struct {
	ch *amqp.Channel
}

func newConsumer(ch *amqp.Channel) *consumer {
	return &consumer{ch: ch}
}

func (c *consumer) Consume(event string) (<-chan []byte, error) {
	messages, err := c.ch.Consume(event, "", autoAck, exclusive, noLocal, noWait, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	resp := make(chan []byte, runtime.NumCPU()*4)
	go func() {
		for m := range messages {
			resp <- m.Body
		}
		close(resp)
	}()

	return resp, nil
}
