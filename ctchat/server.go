package ctchat

import(
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"errors"
	"strconv"
	"gopkg.in/yaml.v2"
)

type(
	Server struct {
		Users map[string]User
		Chatrooms Chatrooms
		WsUpgrader websocket.Upgrader
		HttpDone chan bool
		Port int
		AuthorizedDomains []string
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
	ChatServer.configure()
	ChatServer.startMainChatroom()
	go ChatServer.listenHttp()
	<-ChatServer.HttpDone
}

func(s *Server) configure() error {
	// Temporary structure to load the YAML configuration in
    var config struct {
        Port     string `yaml:"port"`
    	AuthorizedDomains []string `yaml:"authorized_domains"`
    }
    data, err := ioutil.ReadFile("config.yml")
    if err != nil {
        log.Fatal(err)
    }
    if err := yaml.Unmarshal(data, &config); err != nil {
        return err
    }
    port, err := strconv.Atoi(config.Port)
    if err != nil {
        return errors.New("Server config: invalid `port`")
    }
    s.Port = port
    s.AuthorizedDomains = config.AuthorizedDomains
    return nil
}

func (s *Server) startMainChatroom() {
	s.Chatrooms["main"] = &Chatroom{
		Name: "main",
		Users: make(Users, 0),
		Messages: make(Messages, 0),
	}
	s.Chatrooms["main"].StartHub()
	go s.Chatrooms["main"].Hub.Run()
}

func (s *Server) listenHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user User
		
		if !user.HandleAccessControl(w, r) || !user.Authenticate(w, r) {
			return
		}
		user.SendChatroomData(w, s.Chatrooms["main"], http.StatusOK)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if user.HandleAccessControl(w, r) || !user.Authenticate(w, r) {
			return
		}
		user.SendPrivateCommunication(w, "Closing server", http.StatusOK)
		s.HttpDone <- true
	})
	http.Handle("/main", s.Chatrooms["main"])
	if err := http.ListenAndServe(":" + strconv.Itoa(s.Port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}