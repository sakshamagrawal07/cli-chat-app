package main

import (
	"log"

	"github.com/sakshamagrawal07/cli-chat-app.git/server/handler"
	kafkapkg "github.com/sakshamagrawal07/cli-chat-app.git/server/kafka"
	"github.com/sakshamagrawal07/cli-chat-app.git/server/redis"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/utils"
)

func main() {

	err := kafkapkg.InitProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Kafka init error: %v", err)
	}

	port := 8081
	s := utils.NewServer(port)

	// Initialize Redis Pub/Sub
	redis.Init()
	// redis.Subscribe(handler.Broadcast, "chat_room")

	// Register WebSocket route
	router := s.GetRouter()
	router.GET("/ws", handler.HandleWebSocket)

	// Run server
	s.Run()
}
