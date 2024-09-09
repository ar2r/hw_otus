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
	stageOut := stage(stageIn)

	// Передача данных из канала in в канал stageIn
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

	// Обработка данных из канала stageOut для передачи в out
	go func() {
		defer func() {
			close(out)
			// fmt.Println("⛔️ close out channel")
		}()

		for v := range stageOut {
			// Расскомментировать для демонстрации проблемы
			// time.Sleep(time.Millisecond * 100)
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}
