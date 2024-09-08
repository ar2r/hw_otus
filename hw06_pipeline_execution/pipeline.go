package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	wg := sync.WaitGroup{}

	for _, stage := range stages {
		in = runStage(&wg, in, done, stage)
	}

	go func() {
		wg.Wait()
	}()

	return in
}

func runStage(wg *sync.WaitGroup, in In, done In, stage Stage) Out {
	out := make(Bi)
	stageIn := make(Bi)
	wg.Add(2)

	// Перекладывание данных из in в stageIn
	go func() {
		defer func() {
			close(stageIn)
			wg.Done()
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
			wg.Done()
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
