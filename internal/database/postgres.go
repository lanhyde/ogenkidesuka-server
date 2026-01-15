package database

import (
	"database/sql"
	"fmt"

	"github.com/lanhyde/ogenkidesuka-server/internal/config"
)

var DB *sql.DB

func Connect(cfg *config.DatabaseConfig) error {
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db
	fmt.Println("✅Connected to database")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
