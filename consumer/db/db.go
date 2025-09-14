package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var Ctx = context.Background()
var DB *sql.DB

func Init() {
	var err error

	// Open database connection
	// DB, err = sql.Open("postgres", config.PGURL)
	DB, err = sql.Open("postgres", "postgres://postgres:password@localhost:5432/myapp?sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)                 // Maximum number of open connections
	DB.SetMaxIdleConns(5)                  // Maximum number of idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Test the connection with context and timeout
	ctx, cancel := context.WithTimeout(Ctx, 10*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Auto-create messages table if not present
	CreateMessagesTableIfNotExists()
}

// Close gracefully closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
