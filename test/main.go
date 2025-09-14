package main

import (
	"fmt"
	"log"
	"github.com/sakshamagrawal07/cli-chat-app.git/consumer/db"
)

const channel = "whisper:user1:user2"

func main() {
	db.Init()
	defer db.Close()
	messages, err := db.ReadMessages(channel)
	if err != nil {
		log.Fatalf("Error reading messages during test: %v", err)
	}

	for _, msg := range messages {
		fmt.Printf("Message from %s to %s: %s\n", msg.SenderUsername, msg.RecipientUsername, msg.Message)
	}
}