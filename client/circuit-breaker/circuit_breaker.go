package circuitbreaker

import (
	"time"
)

type State int

const (
	StateClosed = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker[T any] struct {
	State         State
	FailureCount  int
	Threshold     int
	Timeout       time.Duration
	LastFailure   time.Time
	DefaultAction func() (T, error)
}

func NewCircuitBreaker[T any](threshold int, timeout time.Duration, defaultAction func() (T, error)) *CircuitBreaker[T] {
	return &CircuitBreaker[T]{
		State:         StateClosed,
		FailureCount:  0,
		Threshold:     threshold,
		Timeout:       timeout,
		LastFailure:   time.Now(),
		DefaultAction: defaultAction,
	}
}

func (cb *CircuitBreaker[T]) Call(action func() (T, error)) (T, error) {
	switch cb.State {
	case StateClosed:
		return cb.stateClosedBehaviour(action)
	case StateOpen:
		return cb.stateOpenBehaviour(action)
	case StateHalfOpen:
		return cb.stateHalfOpenBehaviour(action)
	}

	panic("unknown circuit breaker state")
}

func (cb *CircuitBreaker[T]) stateClosedBehaviour(action func() (T, error)) (T, error) {
	success, err := action()
	if err != nil {
		cb.FailureCount++
		cb.LastFailure = time.Now()

		if cb.FailureCount >= cb.Threshold {
			cb.State = StateOpen

			return cb.DefaultAction()
		}

		return cb.stateClosedBehaviour(action)
	}

	cb.FailureCount = 0

	return success, err
}

func (cb *CircuitBreaker[T]) stateOpenBehaviour(action func() (T, error)) (T, error) {
	if time.Since(cb.LastFailure) >= cb.Timeout {
		cb.State = StateHalfOpen

		return cb.stateHalfOpenBehaviour(action)
	}

	return cb.DefaultAction()
}

func (cb *CircuitBreaker[T]) stateHalfOpenBehaviour(action func() (T, error)) (T, error) {
	success, err := action()
	if err != nil {
		cb.State = StateOpen
		cb.LastFailure = time.Now()

		return cb.DefaultAction()
	}

	cb.State = StateClosed
	cb.FailureCount = 0

	return success, err
}
