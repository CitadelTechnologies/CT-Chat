package ctchat

import(
	"net/http"
	"encoding/json"
	"time"
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
		CreatedAt time.Time `json:"created_at"`
		Chatroom string `json:"chatroom"`
		ExtraData map[string]interface{} `json:"extra_data,omitempty"`
	}
	Messages []Message
)

func (u *User) SendPublicCommunication(w http.ResponseWriter, message string, status int) {
	publicCommunication := PublicCommunication{Message: message}
	sendResponse(w, publicCommunication, status)
}

func (u *User) SendPrivateCommunication(w http.ResponseWriter, message string, status int) {
	privateCommunication := PrivateCommunication{Token: u.Token, Message: message}
	sendResponse(w, privateCommunication, status)

}

func (u *User) SendChatroomData(w http.ResponseWriter, c *Chatroom, status int) {
	chatroomData := ChatroomData{Token: u.Token, Chatroom: c}
	sendResponse(w, chatroomData, status)

}

func sendResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		panic(err)
	}
}

func (u *User) HandleAccessControl(w http.ResponseWriter, r *http.Request) bool {
	origin, isset := r.Header["Origin"]
	if isset != true || len(origin) < 1 || !isAuthorizedDomain(origin[0]) {
		u.SendPublicCommunication(w, "This domain is not authorized", http.StatusForbidden)
		return false
	}
    w.Header().Set("Access-Control-Allow-Origin", origin[0])
    w.Header().Set("Access-Control-Allow-Headers", "accept, authorization")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
    return true
}

func isAuthorizedDomain(origin string) bool {
    for _, domain := range ChatServer.AuthorizedDomains {
        if origin == domain {
            return true
        }
    }
    return false
}