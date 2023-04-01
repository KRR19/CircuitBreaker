package main

import (
	"fmt"
	"net/http"
	"time"
)

var isStoped bool

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		if isStoped {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Result", "Hello")
		writer.WriteHeader(http.StatusOK)
	})
	go func() {
		for range time.Tick(4 * time.Second) {
			isStoped = !isStoped
			fmt.Println("Is server works: ", isStoped)
		}
	}()
	http.ListenAndServe(":8081", nil)
}
