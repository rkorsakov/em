package database

import (
	"fmt"
	"log"
	"os"
	"subscription-service/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dbConfig config.DatabaseConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func Migrate(db *sqlx.DB) error {
	migrationSQL, err := os.ReadFile("internal/migration/001_init.sql")
	if err != nil {
		return fmt.Errorf("error reading migration file: %w", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}

func Close(db *sqlx.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}
