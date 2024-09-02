package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var errCount int
	tasksCh := make(chan Task)    // Канал для задач
	errCh := make(chan error, m)  // Канал для ошибок
	stopCh := make(chan struct{}) // Канал для сигнала остановки

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return
					}
					if err := task(); err != nil {
						errCh <- err
					}
				case <-stopCh:
					return
				}
			}
		}()
	}

	go func() {
		for _, task := range tasks {
			select {
			case tasksCh <- task:
			case <-stopCh:
				return
			}
		}
		close(tasksCh)
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			errCount++
			if errCount >= m {
				close(stopCh)
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}
