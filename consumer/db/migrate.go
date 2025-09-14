package db

import (
	"log"
)

func CreateMessagesTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_username TEXT,
		recipient_username TEXT,
		message TEXT,
		channel TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create messages table: %v", err)
	}
}
