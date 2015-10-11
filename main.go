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
		if r.Method != "POST" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.ParseForm()
		username := r.FormValue("username")
		if username == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Fprintf(w, "Hello " + username)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Closing server ...")
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}