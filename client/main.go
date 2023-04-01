package main

import (
	"fmt"
	circuitBreaker "github.com/KRR19/CircuitBreaker/client/circuit-breaker"
	"net/http"
	"time"
)

func sendRequest() (string, error) {
	res, err := http.Head("http://localhost:8081/hello")
	if err != nil || (res.StatusCode < 200 && res.StatusCode > 299) {
		return "", err
	}

	return res.Header.Get("Result"), nil
}

func main() {
	cb := circuitBreaker.NewCircuitBreaker(5, 2*time.Second)
	for i := 0; i < 100000; i++ {
		result, err := cb.Call(sendRequest)
		if err != nil {
			fmt.Println("Error")
			continue
		}
		fmt.Println(result)
	}
}
