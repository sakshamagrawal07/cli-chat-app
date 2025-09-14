package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

var ctx = context.Background()
var RedisClient *redis.Client
// var channel = "cli-chat-app-messages"

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Publish(msg models.Message) {
	data, _ := json.Marshal(msg)
	if err := RedisClient.Publish(ctx, msg.Channel, data).Err(); err != nil {
		log.Printf("Error publishing message: %v", err)
	}
}

func Subscribe(handler func(models.Message), channel string) {
	pubsub := RedisClient.Subscribe(ctx, channel)
	ch := pubsub.Channel()

	// Start a goroutine to listen for messages
	go func() {
		for msg := range ch {
			var m models.Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err == nil {
				handler(m)
			}
		}
	}()
}