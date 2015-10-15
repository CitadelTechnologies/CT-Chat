package ctchat

type (
	User struct {
		username string
		token string
	}
	Users map[string]User
)
