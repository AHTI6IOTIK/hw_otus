package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	taskChan := make(chan Task, n)
	errChan := make(chan struct{}, n)
	isErrorsExceed := false

	isIgnoreError := m <= 0
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(ctx, isIgnoreError, taskChan, errChan, &wg)
	}

	wg.Add(1)
	go func() {
		step := len(tasks) - 1
		errCount := 0
		for !(errCount > m) {
			select {
			case <-errChan:
				errCount++
				if errCount >= m {
					isErrorsExceed = true
					cancel()
					wg.Done()
					return
				}
			default:
			}

			if step >= 0 {
				taskChan <- tasks[step]
				step--
			}

			if !(step >= 0) && len(taskChan) == 0 {
				wg.Done()
				cancel()
				return
			}
		}
	}()

	wg.Wait()

	close(taskChan)
	close(errChan)

	if isErrorsExceed {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(
	ctx context.Context,
	isIgnoreErrors bool,
	taskChan <-chan Task,
	errChan chan<- struct{},
	wg *sync.WaitGroup,
) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
		}

		select {
		case <-ctx.Done():
			wg.Done()
			return
		case task := <-taskChan:
			err := task()
			if !isIgnoreErrors && err != nil {
				errChan <- struct{}{}
			}
		}
	}
}
