package main

import (
	"log"

	"github.com/sakshamagrawal07/cli-chat-app.git/consumer/db"
	"github.com/sakshamagrawal07/cli-chat-app.git/consumer/kafka"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "chat-messages"
	groupID := "chat-consumer-group"

	db.Init()
	defer db.Close()

	err := kafka.StartKafkaConsumerGroup(brokers, topic, groupID)
	if err != nil {
		log.Fatalf("❌ Consumer failed: %v", err)
	}
}
