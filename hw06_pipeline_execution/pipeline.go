package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func NewDoneStage(done, in In) Out {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case work, ok := <-in:
				if !ok {
					return
				}

				out <- work
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	curr := NewDoneStage(done, in)

	for _, st := range stages {
		curr = NewDoneStage(done, st(curr))
	}

	return curr
}
