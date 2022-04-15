package hw05parallelexecution

import "sync"

type errCounter struct {
	mu      *sync.Mutex
	maxErrs int
	errors  []error
}

func newErrCounter(maxErrs int) *errCounter {
	return &errCounter{
		mu:      &sync.Mutex{},
		maxErrs: maxErrs,
		errors:  make([]error, 0),
	}
}

func (c *errCounter) addErr(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.errors = append(c.errors, err)
}

func (c *errCounter) isExceeded() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxErrs <= 0 {
		return false
	}

	return len(c.errors) >= c.maxErrs
}
