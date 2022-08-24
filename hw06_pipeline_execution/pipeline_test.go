package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestNewDoneStage(t *testing.T) {
	const (
		ms50  = time.Millisecond * 50
		count = 10
	)

	tests := []struct {
		name        string
		doneTimeout time.Duration
		exp         int
	}{
		{
			name:        "exp 0",
			doneTimeout: 0,
		},
		{
			name:        "exp 10",
			doneTimeout: 1 * time.Second,
			exp:         10,
		},
	}

	doneCh := func(timeOut time.Duration) In {
		out := make(chan interface{})
		go func() {
			time.Sleep(timeOut)
			close(out)
		}()
		return out
	}

	in := func() Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for i := 0; i < count; i++ {
				time.Sleep(ms50)
				out <- i
			}
		}()
		return out
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := NewDoneStage(doneCh(tt.doneTimeout), in())

			got := 0
			for {
				_, ok := <-out
				if !ok {
					break
				}
				got++
			}

			require.Equal(t, tt.exp, got)
		})
	}
}
