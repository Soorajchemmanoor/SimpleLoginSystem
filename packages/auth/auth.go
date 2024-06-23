package auth

const USERNAME = "sooraj"
const PASSWORD = "1234"

func AuthLogin(username, password string) bool {
	return username == USERNAME && password == PASSWORD
}
