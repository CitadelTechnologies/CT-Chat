package ctchat

import(
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

type(
	Server struct {
		Users map[string]User
		Chatrooms Chatrooms
		WsUpgrader websocket.Upgrader
		HttpDone chan bool
	}
)

var ChatServer Server

func Start() {
	ChatServer = Server{
		Users: make(map[string]User, 0),
		Chatrooms: make(Chatrooms, 0),
		WsUpgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool{
				return true
			},
		},
		HttpDone: make(chan bool),
	}
	ChatServer.startMainChatroom()
	go ChatServer.listenHttp()
	<-ChatServer.HttpDone
}

func (s *Server) startMainChatroom() {
	s.Chatrooms["main"] = &Chatroom{
		Name: "main",
		Users: make(Users, 0),
		Messages: make(Messages, 0),
		Hub: NewHub(),
	}
	go s.Chatrooms["main"].Hub.Run()
}

func (s *Server) listenHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		s.Chatrooms["main"].AddUser(user)
		user.SendChatroomData(w, s.Chatrooms["main"], http.StatusOK)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if !user.Authenticate(w, r) {
			return
		}
		user.SendPrivateCommunication(w, "Closing server", http.StatusOK)
		s.HttpDone <- true
	})
	http.Handle("/main", s.Chatrooms["main"])
	if err := http.ListenAndServe(":5515", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}