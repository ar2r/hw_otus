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

	// Запустить n горутин
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//  Чтение каналов с задачами и сигналом остановки
			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						// 🤬 Канал задач закрыт
						return
					}
					if err := task(); err != nil {
						// 🚚 Отправить ошибку в канал ошибок
						errCh <- err
					}
				case <-stopCh:
					// 🛑 Получен сигнал остановки
					return
				}
			}
		}()
	}

	// Отправить задачи в канал
	go func() {
		for _, task := range tasks {
			select {
			case tasksCh <- task:
			case <-stopCh:
				// 🔴 Прекратить отправку задач, если получен сигнал остановки
				return
			}
		}
		close(tasksCh)
	}()

	// Ожидание завершения всех горутин
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Ожидание m ошибок
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
