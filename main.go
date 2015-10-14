package main

import(
	"net/http"
	"time"
	"crypto/md5"
	"io"
	"encoding/json"
	"encoding/hex"
	"strings"
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
		var user User
		if r.Method != "POST" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.ParseForm()
		token := getAuthorizationToken(r.Header["Authorization"])
		username := r.FormValue("username")
		if !user.authenticate(w, token, username) {
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

func (u *User) authenticate(w http.ResponseWriter, token string, username string) bool {
	if currentUser, keyExists := users[token]; keyExists {
		u.username = currentUser.username
		u.token = currentUser.token
		return true
	}

	if username == "" {
		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, "You must have an username to sign in")
		return false
	}
	if !isUsernameAvailable(username) {
		w.WriteHeader(http.StatusUnauthorized)
		u.Respond(w, "This username is already taken")
		return false
	}
	newToken := generateToken(username)
	u.username = username
	u.token = newToken
	u.Respond(w, "You are connected")
	users[newToken] = *u
	return false
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

func (u *User) Respond(w http.ResponseWriter, message string) {
	var c Communication
	c.Message = message

	if(u.token != "") {
		c.Token = u.token
	}
	if err := json.NewEncoder(w).Encode(&c); err != nil {
		panic(err)
	}
}

func getAuthorizationToken(authorizationHeader []string) string {
	if len(authorizationHeader) == 0 {
		return ""
	}
	authorization := strings.Split(authorizationHeader[0], " ")
	if len(authorization) > 0 {
		return authorization[1]
	}
	return ""
}