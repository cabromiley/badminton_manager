package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Session store
var Store = sessions.NewCookieStore([]byte("your-secret-key"))

func AuthMiddleware(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// If the user is not authenticated, redirect to the login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}
