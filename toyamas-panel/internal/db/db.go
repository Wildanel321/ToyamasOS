package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(dbPath, defaultAdmin, defaultPass string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQLite pragmas for performance & WAL mode
	_, _ = conn.Exec("PRAGMA journal_mode=WAL;")
	_, _ = conn.Exec("PRAGMA synchronous=NORMAL;")

	database := &DB{Conn: conn}
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := database.ensureAdminUser(defaultAdmin, defaultPass); err != nil {
		log.Printf("[DB WARNING] Failed to seed default admin: %v", err)
	}

	return database, nil
}

func (d *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			action TEXT NOT NULL,
			target TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := d.Conn.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) ensureAdminUser(username, password string) error {
	var count int
	err := d.Conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = d.Conn.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
		if err != nil {
			return err
		}
		log.Printf("[DB SUCCESS] Created default admin user '%s'", username)
	}
	return nil
}

func (d *DB) LogAction(userID int, action, target string) {
	_, _ = d.Conn.Exec("INSERT INTO audit_logs (user_id, action, target, created_at) VALUES (?, ?, ?, ?)",
		userID, action, target, time.Now())
}
