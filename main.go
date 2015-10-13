package main

import(
	"net/http"
	"time"
	"crypto/md5"
	"io"
	"encoding/json"
	"encoding/hex"
)

type (
	User struct {
		username string
		token string
	}
	Users map[string]User
	Communication struct {
		Token string `json:"token"`
		Message string `json:"message"`
	}
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
			c := Communication{Message: "You must have an username to sign in"}
			if err := json.NewEncoder(w).Encode(&c); err != nil {
				panic(err)
			}
			return
		}
		if !isUsernameAvailable(username) {
			w.WriteHeader(http.StatusUnauthorized)
			c := Communication{Message: "This username is already taken"}
			if err := json.NewEncoder(w).Encode(&c); err != nil {
				panic(err)
			}
			return
		}
		token := generateToken(username)
		users[token] = User{username: username, token: token}
		c := Communication{Token: token, Message: "You are connected"}
		if err := json.NewEncoder(w).Encode(&c); err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		c := Communication{Message: "Closing server"}
		if err := json.NewEncoder(w).Encode(&c); err != nil {
			panic(err)
		}
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
	return hex.EncodeToString(hash.Sum(nil))
}
