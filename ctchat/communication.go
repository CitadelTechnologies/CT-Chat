package ctchat

import(
	"net/http"
	"encoding/json"
)

type(
	Communication struct {
		Token string `json:"token"`
		Message string `json:"message"`
	}
)

func (u *User) Respond(w http.ResponseWriter, message string) {
	var c Communication
	c.Message = message

	if(u.token != "") {
		c.Token = u.token
	}
	if err := json.NewEncoder(w).Encode(&c); err != nil {
		panic(err)
	}
}