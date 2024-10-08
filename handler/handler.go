package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/thongsoi/token/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the index page!")
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		token, err := service.AuthenticateUser(db, username, password)
		if err != nil {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}

		// Inject token into HTMX response or set as cookie
		w.Header().Set("HX-Trigger", "logged_in")
		w.Header().Set("HX-Token", token)
	}
}

func ProtectedHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("HX-Token")
		if token == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}

		isValid, err := service.IsValidToken(db, token)
		if err != nil || !isValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed with the protected resource
		w.Write([]byte("This is protected content"))
	}
}
