package config

import (
	"os"
)

type Config struct {
	Port         string
	DBPath       string
	DefaultAdmin string
	DefaultPass  string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./toyamas.db"
	}

	adminUser := os.Getenv("ADMIN_USER")
	if adminUser == "" {
		adminUser = "admin"
	}

	adminPass := os.Getenv("ADMIN_PASS")
	if adminPass == "" {
		adminPass = "toyamas123"
	}

	return &Config{
		Port:         port,
		DBPath:       dbPath,
		DefaultAdmin: adminUser,
		DefaultPass:  adminPass,
	}
}
