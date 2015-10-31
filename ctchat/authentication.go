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
		u.Username = currentUser.Username
		u.Token = currentUser.Token
		return true
	}

	if r.Method == "OPTIONS" {
		u.SendPublicCommunication(w, "", http.StatusOK)
		return false
	}
	if r.Method != "POST" {
		u.SendPublicCommunication(w, "You are not connected", http.StatusUnauthorized)
		return false
	}
	if username == "" {
		u.SendPublicCommunication(w, "You must have an username to sign in", http.StatusForbidden)
		return false
	}
	if !isUsernameAvailable(users, username) {
		u.SendPublicCommunication(w, "This username is already taken", http.StatusUnauthorized)
		return false
	}
	newToken := generateToken(username)
	u.Username = username
	u.Token = newToken
	u.SendPrivateCommunication(w, "You are connected", http.StatusOK)
	ChatServer.Chatrooms["main"].AddUser(*u)
	users[newToken] = *u
	return false
}

func isUsernameAvailable(users map[string]User, username string) bool {
	for _, user := range users {
		if user.Username == username {
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

func (c *Chatroom) AuthenticateConnection(message Message, wsc *WsConnection) bool {
	for _, user := range c.Users {
		if user.Token == message.Token {
			wsc.User = &user
			return true
		}
	}
	return false
}