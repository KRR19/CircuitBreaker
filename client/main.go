package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	circuitBreaker "github.com/KRR19/CircuitBreaker/client/circuit-breaker"
)

func sendRequest() (string, error) {
	res, err := http.Head("http://localhost:8081/hello")
	if err != nil || (res.StatusCode < 200 || res.StatusCode > 299) {
		return "", errors.New("server failed")
	}

	return res.Header.Get("Result"), nil
}

func defaultAction() (string, error) {
	return "DEFAULT", nil
}

func main() {
	cb := circuitBreaker.NewCircuitBreaker(5, 3*time.Second, defaultAction)
	for i := 0; i < 10000000; i++ {
		result, err := cb.Call(sendRequest)
		if err != nil {
			fmt.Println("Error")

			continue
		}
		fmt.Println(result)
	}
}
