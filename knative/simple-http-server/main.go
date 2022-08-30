package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	TIMEOUT_HOUR = 1
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("pending...")
	var count int64
	for {
		count++
		if count == TIMEOUT_HOUR*60*60 {
			break
		}
		if count%60 == 0 {
			fmt.Println(count)
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintf(w, "hello\n")
}

func main() {

	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8093", nil)
}
