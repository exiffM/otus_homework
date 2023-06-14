package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func executeStage(in In, done In, s Stage) Bi {
	out := make(Bi)
	go func() {
		defer close(out)
		out <- s(in)
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	out := make(Bi)

	go func() {
		defer close(out)
		newIn := in
		for _, stage := range stages {
			newIn = executeStage(newIn, done, stage)
		}
		result := <-newIn
		out <- result
	}()
	return out
	// for {
	// 	select {
	// 	case <-done:
	// 		return nil
	// 	default:
	// 		go func() {
	// 			defer close(out)
	// 			newIn := in
	// 			for _, stage := range stages {
	// 				newIn = stage(newIn)
	// 			}
	// 			result := <-newIn
	// 			out <- result
	// 		}()
	// 		return out
	// 	}
	// }

	// go func() {
	// 	for data := range in {
	// 		testChan := make(Bi)
	// 		go func() {
	// 			defer close(testChan)
	// 			testChan <- data
	// 		}()
	// 		nOut := make(Out)
	// 		for _, stage := range stages {
	// 			nOut = stage(testChan)
	// 			testChan = nOut
	// 		}
	// 		result := <-nOut
	// 		out <- result
	// 	}
	// }()
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case data := <-in:
	// 		newIn := make(Bi)
	// 		newIn <- data
	// 		for _, stage := range stages {
	// 			newIn = executeStage(newIn, done, stage)
	// 		}
	// 		result := <-newIn
	// 		close(newIn)
	// 		out <- result
	// 	}
	// }

}
