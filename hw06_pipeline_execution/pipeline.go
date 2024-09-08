package hw06pipelineexecution

import (
	"fmt"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var i int32 = 0
	wg := sync.WaitGroup{}
	for _, stage := range stages {
		fmt.Println("🚚Create stage goroutine", i, stage)
		in = runStage(&wg, i, in, done, stage)
		i++
	}

	go func() {
		wg.Wait()
	}()

	return in
}

func runStage(wg *sync.WaitGroup, i int32, in In, done In, stage Stage) Out {
	wg.Add(1)
	out := make(Bi)
	go func() {
		defer func() {
			fmt.Println(i, "⛔️close outCh")
			wg.Done()
			close(out)
		}()

		stageOut := stage(in)
		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				fmt.Println("\t", i, "🟡new value in stageOut: ", v)
				out <- v
			}
		}
	}()
	return out
}
