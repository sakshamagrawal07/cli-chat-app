package models

import "time"

type Message struct {
	Id                int    `json:"id"`
	SenderUsername    string `json:"sender_username"`
	RecipientUsername string `json:"recipient_username,omitempty"` // Optional for direct messages
	Message           string `json:"message"`
	Channel           string `json:"channel"`
	CreatedAt         time.Time `json:"created_at"` // Use string for JSON serialization
}
