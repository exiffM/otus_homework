package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func executeStage(in In, done In) Bi {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	newIn := in
	for _, stage := range stages {
		stageChan := stage(newIn)
		resultChan := executeStage(stageChan, done)

		newIn = resultChan
	}
	return newIn
}
