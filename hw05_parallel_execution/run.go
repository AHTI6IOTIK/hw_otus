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
	var mu sync.Mutex

	isIgnoreErrors := m <= 0
	step := len(tasks) - 1
	errCount := 0

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for {
				mu.Lock()
				localStep := step
				step--
				localErrCount := errCount
				mu.Unlock()

				if !(localStep >= 0) {
					break
				}

				task := tasks[localStep]

				err := task()
				if err != nil {
					mu.Lock()
					errCount++
					localErrCount = errCount
					mu.Unlock()
				}

				if !isIgnoreErrors && localErrCount >= m {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
