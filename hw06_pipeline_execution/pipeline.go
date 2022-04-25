package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipeIn := make(Bi)
	pipeOut := pipeIn

	for _, s := range stages {
		pipeOut = withDone(done, pipeOut, s)
	}

	go func() {
		defer close(pipeIn)
		for i := range in {
			select {
			case <-done:
				return
			case pipeIn <- i:
			}
		}
	}()

	return pipeOut
}

func withDone(done In, in In, stage Stage) Bi {
	out := make(Bi)
	results := stage(in)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case r, ok := <-results:
				if !ok {
					return // channel is closed
				}

				out <- r
			}
		}
	}()

	return out
}
