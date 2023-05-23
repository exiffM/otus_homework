package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func notifyFinish(qch chan struct{}, n int) {
	for i := 0; i < n; i++ {
		qch <- struct{}{}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	wg := sync.WaitGroup{}
	wg.Add(n)
	defer wg.Wait()
	taskChan := make(chan Task, len(tasks))
	resChan := make(chan error, len(tasks))
	quitChan := make(chan struct{}, n)
	errorFlag := false
	switch m {
	//Случай игнорирования ошибок
	case -1:
		// for i := 0; i < n; i++ {
		// 	go func() {
		// 		defer wg.Done()
		// 		t := <-taskChan
		// 		t()
		// 	}()
		// }
		// for _, task := range tasks {
		// 	taskChan <- task
		// }
		// close(taskChan)

	//Случай когда должно быть 0 ошибок
	case 0:
		// return ErrErrorsLimitExceeded
	//Общий случай при m > 0
	default:
		errorCount := 0
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				for {
					select {
					case <-quitChan:
						fmt.Println("finished")
						return
					case t := <-taskChan:
						resChan <- t()
					}
				}
			}()
		}
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)

		for c := range resChan {
			if c != nil && errorCount < m {
				errorCount++
			} else if errorCount == m {
				errorFlag = true
				notifyFinish(quitChan, n)
				break
			}
		}
	}

	if errorFlag {
		return ErrErrorsLimitExceeded
	} else {
		notifyFinish(quitChan, n)
		return nil
	}
}
