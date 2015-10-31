package ctchat

import(
	"net/http"
	"encoding/json"
)

type(
	PublicCommunication struct {
		Message string `json:"message"`
	}
	PrivateCommunication struct {
		Token string `json:"token"`
		Message string `json:"message"`
	}
	ChatroomData struct {
		Token string `json:"token"`
		Chatroom *Chatroom `json:"chatroom"`
	}
	Message struct {
		Token string `json:"token,omitempty"`
		Type string `json:"type"`
		Author string `json:"author"`
		Content string `json:"content"`
		Chatroom string `json:"chatroom"`
		ExtraData map[string]interface{} `json:"extra_data"`
	}
	Messages []Message
)

func (u *User) SendPublicCommunication(w http.ResponseWriter, message string, status int) {
	publicCommunication := PublicCommunication{Message: message}

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "accept, authorization")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&publicCommunication); err != nil {
		panic(err)
	}
}

func (u *User) SendPrivateCommunication(w http.ResponseWriter, message string, status int) {
	privateCommunication := PrivateCommunication{Token: u.Token, Message: message}

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "accept, authorization")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&privateCommunication); err != nil {
		panic(err)
	}
}

func (u *User) SendChatroomData(w http.ResponseWriter, c *Chatroom, status int) {
	chatroomData := ChatroomData{Token: u.Token, Chatroom: c}

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "accept, authorization")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&chatroomData); err != nil {
		panic(err)
	}
}