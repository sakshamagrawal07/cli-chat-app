package config

import "os"

var (
	KafkaBroker = os.Getenv("KAFKA_BROKER")          // e.g. "localhost:9092"
	KafkaTopic  = os.Getenv("KAFKA_TOPIC")           // e.g. "messages"
	PGURL       = os.Getenv("POSTGRES_URL")          // e.g. "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
)