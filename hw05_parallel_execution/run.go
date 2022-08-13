package hw05parallelexecution

import (
	"runtime"
	"sync"
)

type Task func() error

func worker(wg *sync.WaitGroup, ch <-chan Task, counter *Counter) {
	defer wg.Done()

	for task := range ch {
		if ok := counter.IsMaximumAttemptsValue(); ok {
			return
		}

		if err := task(); err != nil {
			if ok := counter.Increment(); !ok {
				return
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error { // n - numbersOfGoroutine, m-availableErrorsBeforeStop
	ch := make(chan Task, runtime.NumCPU())
	counter := NewCounter(int32(m))

	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(wg, ch, counter)
	}

	var failed bool
	for _, task := range tasks {
		if counter.IsMaximumAttemptsValue() {
			failed = true
			break
		}

		ch <- task
	}

	close(ch)
	wg.Wait()

	var err error
	if failed || counter.IsMaximumAttemptsValue() {
		err = ErrErrorsLimitExceeded
	}

	return err
}
