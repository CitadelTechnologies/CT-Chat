package ctchat

import(
	"net/http"
	"time"
	"crypto/md5"
	"io"
	"encoding/hex"
	"strings"
)

func (u *User) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	token := getAuthorizationToken(r.Header["Authorization"])
	r.ParseForm()
	username := r.FormValue("username")

	users := ChatServer.Users

	if currentUser, keyExists := users[token]; keyExists {
		u.username = currentUser.username
		u.token = currentUser.token
		return true
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if username == "" {
		w.WriteHeader(http.StatusForbidden)
		u.SendPublicCommunication(w, "You must have an username to sign in")
		return false
	}
	if !isUsernameAvailable(users, username) {
		w.WriteHeader(http.StatusUnauthorized)
		u.SendPublicCommunication(w, "This username is already taken")
		return false
	}
	newToken := generateToken(username)
	u.username = username
	u.token = newToken
	u.SendPrivateCommunication(w, "You are connected")
	users[newToken] = *u
	return false
}

func isUsernameAvailable(users Users, username string) bool {
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