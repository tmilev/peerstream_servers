package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/backend", requestBackend)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Fatal error: failed to run http server. %v", err)
	} else {
		fmt.Printf("Server done. ")
	}
}

func requestBackend(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello world!")
	if err != nil {
		fmt.Printf("Error writing response. ")
	}
}
