package main

import(
	"fmt"
	"net/http"
)

func main() {

	httpDone := make(chan bool)
	go listenHttp(httpDone)
	<-httpDone
}

func listenHttp(done chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Closing server ...")
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}