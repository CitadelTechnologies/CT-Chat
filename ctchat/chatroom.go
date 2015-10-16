package ctchat

import(
	"time"
)

type(
	Chatroom struct {
		Name string `json:"name"`
		Users []User `json:"users"`
		Messages Messages `json:"messages"`
	}
	Chatrooms map[string]Chatroom
	Message struct {
		Author *User `json:"author"`
		CreatedAt time.Time `json:"created_at"`
		Content string `json:"content"`
	}
	Messages []Message
)