package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

var (
	producer sarama.SyncProducer
	topic    = "chat-messages"
)

func InitProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1
	config.Version = sarama.V2_6_0_0

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("Failed to start Kafka producer: %v\n", err)
		return err
	}
	log.Println("Kafka producer initialized")
	return nil
}

func ProduceMessage(msg models.Message) error {
	log.Println("[DEBUG] Producing message:", msg)
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(payload),
		Timestamp: time.Now(),
	}

	_, _, err = producer.SendMessage(kafkaMsg)
	return err
}
