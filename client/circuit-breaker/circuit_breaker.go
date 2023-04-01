package circuitBreaker

import (
	"errors"
	"time"
)

type CircuitBreakerState int

const (
	StateClosed = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	State         CircuitBreakerState
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
		success, err := action()
		if err != nil {
			cb.FailureCount++
			cb.LastFailure = time.Now()
			if cb.FailureCount >= cb.Threshold {
				cb.State = StateOpen
			}
		} else {
			cb.FailureCount = 0
		}
		return success, err

	case StateOpen:
		if time.Since(cb.LastFailure) >= cb.Timeout {
			cb.State = StateHalfOpen
			return cb.Call(action)
		}
		return cb.DefaultAction()

	case StateHalfOpen:
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

	return "", errors.New("unknown circuit breaker state")
}
