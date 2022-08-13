package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Counter struct {
	maxValue int32
	mu       *sync.Mutex
	count    int32
}

func NewCounter(maxValue int32) *Counter {
	return &Counter{
		maxValue: maxValue,
		mu:       &sync.Mutex{},
	}
}

func (c Counter) isMaximumAttemptsValue() bool {
	return c.count >= c.maxValue && c.maxValue != 0
}

func (c *Counter) IsMaximumAttemptsValue() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.isMaximumAttemptsValue()
}

func (c *Counter) Increment() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if isDone := c.isMaximumAttemptsValue(); isDone {
		return false
	}

	c.count++

	return true
}
