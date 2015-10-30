package ctchat

import(
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"html"
)

type(
	Chatroom struct {
		Name string `json:"name"`
		Users []User `json:"users"`
		Messages Messages `json:"messages"`
		Hub *WsHub `json:"-"`
	}
	Chatrooms map[string]*Chatroom
	WsConnection struct {
		// The websocket connection.
		Conn *websocket.Conn
		// Buffered channel of outbound messages.
		Send chan []byte
		// The hub.
		Hub *WsHub
	}
	WsHub struct {
		// Registered connections.
		Connections map[*WsConnection]bool
		// Inbound messages from the connections.
		Broadcast chan []byte
		// Register requests from the connections.
		Register chan *WsConnection
		// Unregister requests from connections.
		Unregister chan *WsConnection
	}
)

func (chatroom *Chatroom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := ChatServer.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	c := &WsConnection{
		Send: make(chan []byte, 256),
		Conn: ws,
		Hub: chatroom.Hub,
	}
	c.Hub.Register <- c
	defer func() { c.Hub.Unregister <- c }()
	go c.Write()
	c.Read(chatroom)
}

func NewHub() *WsHub {
	return &WsHub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *WsConnection),
		Unregister:  make(chan *WsConnection),
		Connections: make(map[*WsConnection]bool),
	}
}

func (hub *WsHub) Run() {
	for {
		select {
		case c := <-hub.Register:
			hub.Connections[c] = true
		case c := <-hub.Unregister:
			if _, ok := hub.Connections[c]; ok {
				delete(hub.Connections, c)
				close(c.Send)
			}
		case m := <-hub.Broadcast:
			for c := range hub.Connections {
				select {
				case c.Send <- m:
				default:
					delete(hub.Connections, c)
					close(c.Send)
				}
			}
		}
	}
}

func (wsc *WsConnection) Read(chatroom *Chatroom) {
	for {
		_, data, err := wsc.Conn.ReadMessage()
		if err != nil {
			break
		}
		jsonBroadcast, err := json.Marshal(chatroom.AddMessage(data))
		if err != nil {
			panic(err)
		}
		wsc.Hub.Broadcast <- jsonBroadcast
	}
	wsc.Conn.Close()
}

func (c *Chatroom) AddMessage(data []byte) Message {
	message := Message{}
	json.Unmarshal([]byte(data), &message)
	message.Token = ""
	message.Type = "message"
	message.Content = html.EscapeString(strings.TrimSpace(message.Content))
	c.Messages = append(c.Messages, message)
	if len(c.Messages) > 100 {
		// Remove the older message when the limit is reached
		c.Messages = append(c.Messages[:0], c.Messages[1:]...)
	}
	return message
}

func (wsc *WsConnection) Write() {
	for message := range wsc.Send {
		err := wsc.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	wsc.Conn.Close()
}

func (c *Chatroom) AddUser(data ...User) {
    m := len(c.Users)
    n := m + len(data)
    if n > cap(c.Users) {
        newSlice := make(Users, (n+1)*2)
        copy(newSlice, c.Users)
        c.Users = newSlice
    }
    c.Users = c.Users[0:n]
    copy(c.Users[m:n], data)
}