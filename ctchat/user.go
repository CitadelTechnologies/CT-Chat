package ctchat

type (
	User struct {
		Username string `json:"username"`
		Token string `json:"-"`
	}
	Users []User
)