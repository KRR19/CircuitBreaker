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
	State        CircuitBreakerState
	FailureCount int
	Threshold    int
	Timeout      time.Duration
	LastFailure  time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		State:        StateClosed,
		FailureCount: 0,
		Threshold:    threshold,
		Timeout:      timeout,
		LastFailure:  time.Now(),
	}
}

func (cb *CircuitBreaker) Call(protectedFunc func() (string, error)) (string, error) {
	switch cb.State {
	case StateClosed:
		success, err := protectedFunc()
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
			return cb.Call(protectedFunc)
		}
		return "", errors.New("circuit breaker is open")

	case StateHalfOpen:
		success, err := protectedFunc()
		if err != nil {
			cb.State = StateOpen
			cb.LastFailure = time.Now()
		} else {
			cb.State = StateClosed
			cb.FailureCount = 0
		}
		return success, err
	}

	return "", errors.New("unknown circuit breaker state")
}
