package main

import(
	"fmt"
	"net/http"
	"time"
	"crypto/md5"
	"io"
)

type (
	User struct {
		username string
		token string
	}
	Users map[string]User
)

var users Users

func main() {
	users = make(Users, 0)
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
			fmt.Fprintf(w, "You must have an username to sign in")
			return
		}
		if !isUsernameAvailable(username) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "This username is already taken")
			return
		}
		token := generateToken(username)
		users[token] = User{username: username, token: token}
		fmt.Fprintf(w, "Hello " + username)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Closing server ...")
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}

func isUsernameAvailable(username string) bool {
	for _, user := range users {
		if(user.username == username) {
			return false
		}
	}
	return true
}

func generateToken(username string) string {
	hash := md5.New()
	io.WriteString(hash, username)
	io.WriteString(hash, time.Now().Format(time.UnixDate))
	return string(hash.Sum(nil))
}