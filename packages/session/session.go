// session/session.go
package session

import (
	"net/http"
	"time"
)

const sessionCookieName = "session"

func SetSession(w http.ResponseWriter, username string) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    username,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GetSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
