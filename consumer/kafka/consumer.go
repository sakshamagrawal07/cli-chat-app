package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
	"github.com/sakshamagrawal07/cli-chat-app.git/consumer/db"
)

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for data := range claim.Messages() {
		log.Printf("📥 Message received | Topic: %s | Partition: %d | Offset: %d | Key: %s | Value: %s\n", data.Topic, data.Partition, data.Offset, string(data.Key), string(data.Value))

		var msg models.Message
		if err := json.Unmarshal(data.Value, &msg); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			return err
		}
		// Process the message (e.g., save to database)
		if err := db.InsertMessage(msg); err != nil {
			log.Println("Error saving message to database:", err)
			return err
		}

		// Mark as processed
		session.MarkMessage(data, "")
	}
	return nil
}

func StartKafkaConsumerGroup(brokers []string, topic string, groupID string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Handle Ctrl+C
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Println("🔴 Interrupt detected, stopping consumer...")
		cancel()
	}()

	log.Println("🚀 Starting Kafka consumer group...")

	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, ConsumerGroupHandler{}); err != nil {
			log.Printf("⚠️ Error during consumption: %v", err)
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
