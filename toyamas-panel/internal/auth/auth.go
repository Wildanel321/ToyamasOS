package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"toyamas-panel/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func Authenticate(database *db.DB, username, password string) (*User, error) {
	var user User
	var hash string

	err := database.Conn.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

func CreateSession(database *db.DB, userID int) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := database.Conn.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateSession(database *db.DB, token string) (*User, error) {
	if token == "" {
		return nil, errors.New("empty session token")
	}

	var user User
	var expiresAt time.Time

	err := database.Conn.QueryRow(`
		SELECT u.id, u.username, s.expires_at 
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.token = ?`, token).Scan(&user.ID, &user.Username, &expiresAt)

	if err != nil {
		return nil, errors.New("invalid or expired session")
	}

	if time.Now().After(expiresAt) {
		_ = DeleteSession(database, token)
		return nil, errors.New("session expired")
	}

	return &user, nil
}

func DeleteSession(database *db.DB, token string) error {
	_, err := database.Conn.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}
