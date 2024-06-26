package hw05parallelexecution

import (
	"errors"
	"math"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// precheck
	if len(tasks) == 0 || n <= 0 {
		return nil
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	// init
	handlersCount := int(math.Min(float64(n), float64(len(tasks))))
	rets := make(chan error, handlersCount)
	taskQueue := make(chan Task, handlersCount)
	wg := new(sync.WaitGroup)
	defer func() {
		close(taskQueue)
		wg.Wait()
		close(rets)
	}()
	// handlers
	for i := 0; i < handlersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskQueue {
				rets <- t()
			}
		}()
	}
	// work
	for _, t := range tasks[:handlersCount] {
		taskQueue <- t
	}
	enqueued := handlersCount
	for i := 0; i < len(tasks); i++ {
		if ret := <-rets; ret != nil {
			if m--; m == 0 {
				return ErrErrorsLimitExceeded
			}
		}
		if enqueued < len(tasks) {
			taskQueue <- tasks[enqueued]
			enqueued++
		}
	}
	return nil
}
