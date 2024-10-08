package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"your_project/repository"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(db *sql.DB, username, password string) (string, error) {
	user, err := repository.GetUserByUsername(db, username)
	if err != nil || user == nil {
		return "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	err = repository.StoreToken(db, token, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsValidToken(db *sql.DB, token string) (bool, error) {
	_, err := repository.ValidateToken(db, token)
	return err == nil, err
}
