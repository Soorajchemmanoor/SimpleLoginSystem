package validation

func ValidateForm(username, password string) (string, string) {
	uErr, pErr := "", ""

	if username == "" {
		uErr = "Username is required"
	}
	if password == "" {
		pErr = "Password is required"
	}
	return uErr, pErr
}
