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
		Chatroom Chatroom `json:"chatroom"`
	}
)

func (u *User) SendPublicCommunication(w http.ResponseWriter, message string) {
	publicCommunication := PublicCommunication{Message: message}

	if err := json.NewEncoder(w).Encode(&publicCommunication); err != nil {
		panic(err)
	}
}

func (u *User) SendPrivateCommunication(w http.ResponseWriter, message string) {
	privateCommunication := PrivateCommunication{Token: u.token, Message: message}

	if err := json.NewEncoder(w).Encode(&privateCommunication); err != nil {
		panic(err)
	}
}

func (u *User) SendChatroomData(w http.ResponseWriter, c Chatroom) {
	chatroomData := ChatroomData{Token: u.token, Chatroom: c}

	if err := json.NewEncoder(w).Encode(&chatroomData); err != nil {
		panic(err)
	}
}