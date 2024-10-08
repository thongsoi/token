package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/thongsoi/gorilla-sessions/models"
)

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password FROM users WHERE username = $1"
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func StoreToken(db *sql.DB, token string, userID int) error {
	query := "INSERT INTO tokens (token, user_id) VALUES ($1, $2)"
	_, err := db.Exec(query, token, userID)
	return err
}

func ValidateToken(db *sql.DB, token string) (int, error) {
	var userID int
	query := "SELECT user_id FROM tokens WHERE token = $1"
	err := db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
