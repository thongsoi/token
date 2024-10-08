package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/thongsoi/token/service"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println("Attempting login with username:", username) // Log username

		// Try to authenticate the user
		token, err := service.AuthenticateUser(db, username, password)
		if err != nil {
			fmt.Println("Login error:", err) // Log the error
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}

		fmt.Println("Login successful, generated token:", token) // Log token

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
