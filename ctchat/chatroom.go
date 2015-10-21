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


func (c Chatroom) AddUser(data ...User) {
    m := len(c.Users)
    n := m + len(data)
    if n > cap(c.Users) {
        newSlice := make([]User, (n+1)*2)
        copy(newSlice, c.Users)
        c.Users = newSlice
    }
    c.Users = c.Users[0:n]
    copy(c.Users[m:n], data)
}