package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.

	// Первый вариант
	// out := make(Bi, len(stages))
	// go func() {
	// 	defer close(out)
	// 	for _, stage := range stages {
	// 		out <- stage(in)
	// 	}
	// }()

	// Второй вариант
	// out := make(Bi, len(stages))
	// stageIn := make(Bi)
	// go func() {
	// 	defer close(stageIn)
	// 	for data := range in {
	// 		stageIn <- data
	// 	}
	// }()

	// go func() {
	// 	defer close(out)
	// 	for _, stage := range stages {
	// 		out <- stage(stageIn)
	// 	}
	// }()

	// Третий вариант
	// out := make(Bi, len(stages))
	// stageIn := make(Bi)
	// go func() {
	// 	defer close(stageIn)
	// 	for data := range in {
	// 		stageIn <- data
	// 	}
	// }()

	//
	// 	defer close(out)
	// 	for _, stage := range stages {
	// 		out <- stage(stageIn)
	// 	}

	return out
}
