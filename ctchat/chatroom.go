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
		User *User
		// The websocket connection.
		Conn *websocket.Conn
		// Buffered channel of outbound messages.
		Send chan []byte
		// The hub.
		Hub *WsHub
	}
	WsHub struct {
		Chatroom *Chatroom
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

func (c *Chatroom) StartHub() {
	c.Hub = &WsHub{
		Chatroom:    c,
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
				if c.User != nil {
					hub.disconnectUser(c)
				}
				delete(hub.Connections, c)
				close(c.Send)
			}
		case m := <-hub.Broadcast:
			hub.broadcast(m)
		}
	}
}

func (hub *WsHub) disconnectUser(c *WsConnection) {
	jsonBroadcast, err := json.Marshal(Message{
		Type: "notification",
		Author: c.User.Username,
		Content: c.User.Username +  " is disconnected" ,
		Chatroom: hub.Chatroom.Name,
		ExtraData: map[string]interface{}{"notification_type": "disconnection"},
	})
	if err != nil {
		panic(err)
	}
	hub.broadcast(jsonBroadcast)
	hub.Chatroom.DeleteUser(c.User)
	c.User = nil
}

func (hub *WsHub) broadcast(m []byte) {
	for c := range hub.Connections {
		select {
		case c.Send <- m:
		default:
			if c.User != nil {
				hub.disconnectUser(c)
			}
			delete(hub.Connections, c)
			close(c.Send)
		}
	}
}

func (wsc *WsConnection) Read(chatroom *Chatroom) {
	for {
		_, data, err := wsc.Conn.ReadMessage()
		if err != nil {
			break
		}
		message := Message{}
		json.Unmarshal([]byte(data), &message)
		switch {
			case message.Type == "authentication":
				if chatroom.AuthenticateConnection(message, wsc) == true {
					message.Type = "notification"
					message.Content = message.Author + " is connected"
					message.ExtraData = map[string]interface{}{"notification_type": "connection"}
					jsonBroadcast, err := json.Marshal(message)
					if err != nil {
						panic(err)
					}
					wsc.Hub.Broadcast <- jsonBroadcast
				}
			default:
				message.Type = "message"
				jsonBroadcast, err := json.Marshal(chatroom.AddMessage(message))
				if err != nil {
					panic(err)
				}
				wsc.Hub.Broadcast <- jsonBroadcast
		}
	}
	wsc.Conn.Close()
}

func (c *Chatroom) AddMessage(message Message) Message {
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

func (c *Chatroom) DeleteUser(u *User) {
	for index, user := range c.Users {
		if user.Token == u.Token {
			c.Users = append(c.Users[:index], c.Users[index+1:]...)
			delete(ChatServer.Users, u.Token)
		}
	}
}