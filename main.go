package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thongsoi/token/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables or set up database connection details
	connStr := "postgres://dev1:dev1pg@localhost:5432/biomassx?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to ping the database:", err)
	}

	// Create a new router
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/login", handler.LoginHandler(db)).Methods("POST")
	router.HandleFunc("/protected", handler.ProtectedHandler(db)).Methods("GET")

	// Middleware to log incoming requests
	router.Use(loggingMiddleware)

	// Serve static files (if any, like HTML/CSS/JS)
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// loggingMiddleware logs incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
