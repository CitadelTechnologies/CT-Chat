package ctchat

import(
	"net/http"
)

type(
	Server struct {
		Users map[string]User
		Chatrooms Chatrooms
	}
)

var ChatServer Server

func Start() {
	ChatServer = Server{
		Users: make(map[string]User, 0),
		Chatrooms: make(Chatrooms, 0),
	}
	ChatServer.startMainChatroom()
	httpDone := make(chan bool)
	go listenHttp(httpDone)
	<-httpDone
}

func (s *Server) startMainChatroom() {
	s.Chatrooms["main"] = Chatroom{Name: "main", Users: make([]User, 0), Messages: make(Messages, 0)}
}

func listenHttp(done chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		ChatServer.Chatrooms["main"].Users[len(ChatServer.Chatrooms["main"].Users) + 1] = user
		user.SendChatroomData(w, ChatServer.Chatrooms["main"], http.StatusOK)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		user.SendPrivateCommunication(w, "Closing server", http.StatusOK)
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}