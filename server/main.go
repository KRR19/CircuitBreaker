package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var isStopped bool

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		if isStopped {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}
		writer.Header().Add("Result", "Hello")
		writer.WriteHeader(http.StatusOK)
	})
	go func() {
		for range time.Tick(4 * time.Second) {
			isStopped = !isStopped
			fmt.Println("Is server works: ", isStopped)
		}
	}()
	log.Fatal(http.ListenAndServe(":8081", nil))
}
