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
	taskChan := make(chan Task)
	errChan := make(chan struct{}, n+m)

	isErrorsExceed := false
	isIgnoreError := m <= 0
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(isIgnoreError, taskChan, errChan, &wg)
	}

	errCnt := 0
	for _, t := range tasks {
		select {
		case <-errChan:
			errCnt++
		default:
		}

		if !isIgnoreError && errCnt >= m {
			isErrorsExceed = true
			break
		}

		taskChan <- t
	}
	close(taskChan)
	wg.Wait()

	if isErrorsExceed {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(isIgnoreErrors bool, taskChan <-chan Task, errChan chan<- struct{}, wg *sync.WaitGroup) {
	for task := range taskChan {
		err := task()
		if !isIgnoreErrors && err != nil {
			errChan <- struct{}{}
		}
	}

	wg.Done()
}
