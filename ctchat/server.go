package ctchat

import(
	"net/http"
)

type(
	Server struct {
		Users Users
		Chatrooms Chatrooms
	}
)

var ChatServer Server

func Start() {
	ChatServer = Server{
		Users: make(Users, 0),
		Chatrooms: make(Chatrooms, 0),
	}
	ChatServer.startMainChatroom()
	httpDone := make(chan bool)
	go listenHttp(httpDone)
	<-httpDone
}

func (s *Server) startMainChatroom() {
	s.Chatrooms["main"] = Chatroom{Name: "main", Users: make(Users, 0), Messages: make(Messages, 0)}
}

func listenHttp(done chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		user.SendChatroomData(w, ChatServer.Chatrooms["main"])
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		user.SendPrivateCommunication(w, "Closing server")
		done <- true
	})
	http.ListenAndServe(":5515", nil)
}