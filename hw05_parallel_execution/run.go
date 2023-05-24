package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type safeCounter struct {
	mutex   sync.RWMutex
	counter int
}

func (sc *safeCounter) Increment() {
	sc.mutex.Lock()
	sc.counter++
	sc.mutex.Unlock()
}

func (sc *safeCounter) Counter() int {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.counter
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)
	taskChan := make(chan Task, len(tasks))
	errorCounter := safeCounter{}
	switch m {
	// Случай игнорирования ошибок
	case -1:
		executor := func() {
			defer wg.Done()
			for task := range taskChan {
				task()
			}
		}
		for i := 0; i < n; i++ {
			go executor()
		}
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
		wg.Wait()
	// Случай когда должно быть 0 ошибок
	case 0:
		return ErrErrorsLimitExceeded
	// Общий случай при m > 0
	default:
		executor := func() {
			defer wg.Done()
			for task := range taskChan {
				if err := task(); err != nil && errorCounter.Counter() < m {
					errorCounter.Increment()
				} else if errorCounter.Counter() == m {
					return
				}
			}
		}
		for i := 0; i < n; i++ {
			go executor()
		}
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
		wg.Wait()
		if errorCounter.counter == m {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
