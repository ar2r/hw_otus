package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = runStage(in, done, stage)
	}

	return in
}

func runStage(in In, done In, stage Stage) Out {
	out := make(Bi)
	stageIn := make(Bi)

	// Перекладывание данных из in в stageIn
	go func() {
		defer func() {
			close(stageIn)
			// fmt.Println("⛔️ close stageIn channel")
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case stageIn <- v:
				}
			}
		}
	}()

	// Обработка данных из канала stageIn
	go func() {
		defer func() {
			close(out)
			// fmt.Println("⛔️ close out channel")
		}()

		for v := range stage(stageIn) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
