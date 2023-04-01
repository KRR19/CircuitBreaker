package circuitBreaker

import (
	"time"
)

type State int

const (
	StateClosed = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	State         State
	FailureCount  int
	Threshold     int
	Timeout       time.Duration
	LastFailure   time.Time
	DefaultAction func() (string, error)
}

func NewCircuitBreaker(threshold int, timeout time.Duration, defaultAction func() (string, error)) *CircuitBreaker {
	return &CircuitBreaker{
		State:         StateClosed,
		FailureCount:  0,
		Threshold:     threshold,
		Timeout:       timeout,
		LastFailure:   time.Now(),
		DefaultAction: defaultAction,
	}
}

func (cb *CircuitBreaker) Call(action func() (string, error)) (string, error) {
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

func (cb *CircuitBreaker) stateClosedBehaviour(action func() (string, error)) (string, error) {
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

func (cb *CircuitBreaker) stateOpenBehaviour(action func() (string, error)) (string, error) {
	if time.Since(cb.LastFailure) >= cb.Timeout {
		cb.State = StateHalfOpen
		return cb.stateHalfOpenBehaviour(action)
	}

	return cb.DefaultAction()
}

func (cb *CircuitBreaker) stateHalfOpenBehaviour(action func() (string, error)) (string, error) {
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
