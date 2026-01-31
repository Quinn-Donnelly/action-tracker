package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Connect establishes a database connection with the given connection string
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations if needed
func Migrate(db *sql.DB) error {
	// The schema is created by init-db.sql when the container starts
	// This function can be used for additional migrations if needed
	return nil
}
