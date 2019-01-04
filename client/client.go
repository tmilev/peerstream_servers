package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/request", requestBackend)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Printf(
		"Proceeding to listen to port 3000. Please  open \nhttp://localhost:3000\n in your web browser. \n",
	)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Fatal error: failed to run http server. %v", err)
	} else {
		fmt.Printf("Server done. ")
	}
}