package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task, n+m)
	wg := &sync.WaitGroup{}
	counter := newErrCounter(m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go consumer(wg, counter, taskCh)
	}

	wg.Add(1)
	go producer(counter, wg, tasks, taskCh)

	wg.Wait()

	if counter.isExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func producer(counter *errCounter, wg *sync.WaitGroup, tasks []Task, taskCh chan Task) {
	defer wg.Done()
	defer close(taskCh)

	for _, task := range tasks {
		if counter.isExceeded() {
			return
		}

		taskCh <- task
	}
}

func consumer(wg *sync.WaitGroup, counter *errCounter, taskCh chan Task) {
	defer wg.Done()

	for task := range taskCh {
		if err := task(); err != nil {
			counter.addErr(err)
		}

		if counter.isExceeded() {
			return
		}
	}
}
