package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	port = 8080
)

var (
	db  *sql.DB
	key = []byte("your-secret-key") // Replace this with a secure secret key
)

func main() {
	var err error

	// Connect to the database
	db, err = connectDatabase()
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer db.Close()

	// Set up the routes
	mux := http.NewServeMux()
	mux = setupRoutes(mux)

	// Start the server
	fmt.Printf("Listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func setupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/api/auth/pin-login", pinLoginHandler)
	mux.HandleFunc("/api/canary/checkin", canaryCheckinHandler)
	mux.HandleFunc("/api/canary/status", canaryStatusHandler)
	mux.HandleFunc("/api/events", getUserEventsHandler)
	mux.HandleFunc("/api/events", createEventHandler)
	mux.HandleFunc("/api/events/", updateEventHandler)
	mux.HandleFunc("/api/events/", deleteEventHandler)
	mux.HandleFunc("/api/messages", getMessagesHandler)
	mux.HandleFunc("/api/messages", sendMessageHandler)
	mux.HandleFunc("/api/messages/", getMessageDetailsHandler)
	mux.HandleFunc("/api/messages/", updateMessageHandler)
	mux.HandleFunc("/api/messages/", deleteMessageHandler)
	mux.HandleFunc("/api/users", getUsersHandler)
	mux.HandleFunc("/api/users/", getUserDetailsHandler)
	mux.HandleFunc("/api/users/", updateUserHandler)
	mux.HandleFunc("/api/users/", deleteUserHandler)

	return mux
}

func connectDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./fell_es.db")
	if err != nil {
		return nil, err
	}

	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		pin TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		canary_last_checked DATETIME
	);
	`

	_, err = db.Exec(usersTable)
	if err != nil {
		return nil, err
	}

	// Create events table
	eventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		start DATETIME NOT NULL,
		end DATETIME NOT NULL,
		location TEXT NOT NULL,
		organizer_id INTEGER NOT NULL,
		FOREIGN KEY (organizer_id) REFERENCES users (id)
	);
	`

	_, err = db.Exec(eventsTable)
	if err != nil {
		return nil, err
	}

	// Create messages table
	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		sent_at DATETIME NOT NULL,
		FOREIGN KEY (sender_id) REFERENCES users (id),
		FOREIGN KEY (receiver_id) REFERENCES users (id)
	);
	`

	_, err = db.Exec(messagesTable)
	if err != nil {
		return nil, err
	}

	// Create users_events table (for managing event participants)
	usersEventsTable := `
	CREATE TABLE IF NOT EXISTS users_events (
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		PRIMARY KEY (user_id, event_id),
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (event_id) REFERENCES events (id)
	);
	`

	_, err = db.Exec(usersEventsTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupJWT() (*jwt.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token, nil
}
