package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func GenerateOAuthState(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(10 * time.Minute),
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	return state
}

func ValidateOAuthState(r *http.Request, stateParam string) bool {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		return false
	}
	return stateParam == stateCookie.Value
}

func ClearOAuthState(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}
