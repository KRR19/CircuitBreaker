package main

import (
	"fmt"
	"net/http"
)

func sendRequest() (string, error) {
	res, err := http.Head("http://localhost:8081/hello")
	if err != nil || (res.StatusCode < 200 && res.StatusCode > 299) {
		return "", err
	}

	return res.Header.Get("Result"), nil
}

func main() {
	for i := 0; i < 10000000; i++ {
		result, err := sendRequest()
		if err != nil {
			fmt.Println("Error")
			continue
		}
		fmt.Println(result)
	}
}
