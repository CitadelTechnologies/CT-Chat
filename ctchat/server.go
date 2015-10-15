package ctchat

import(
	"net/http"
)

type(
	Server struct {
		Users Users
	}
)

var ChatServer Server

func Start() {
	ChatServer = Server{Users: make(Users, 0)}
	httpDone := make(chan bool)
	go listenHttp(httpDone)
	<-httpDone
}

func listenHttp(done chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		user.Respond(w, "Connected")
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var user User
		w.WriteHeader(http.StatusUnauthorized)
		user.Respond(w, "Closing server")
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}