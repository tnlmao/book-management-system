package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Kafka topic for book events
const BookEventsTopic = "book_events"

// BookEvent represents the event structure
type BookEvent struct {
	Event string      `json:"event"` // "create", "update", "delete"
	Book  interface{} `json:"book"`
	Time  string      `json:"time"`
}

// KafkaProducer structure
type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

// NewKafkaProducer initializes a Kafka producer
func NewKafkaProducer(broker string) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	return &KafkaProducer{
		Producer: producer,
		Topic:    BookEventsTopic,
	}
}

// PublishBookEvent sends an event to Kafka
func (p *KafkaProducer) PublishBookEvent(eventType string, book interface{}) {
	event := BookEvent{
		Event: eventType,
		Book:  book,
		Time:  time.Now().Format(time.RFC3339),
	}

	// Convert event to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal book event: %v", err)
		return
	}

	// Send message to Kafka
	err = p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)

	if err != nil {
		log.Printf("Failed to publish %s event: %v", eventType, err)
	} else {
		log.Printf("%s event published for book: %+v", eventType, book)
	}
}

// Close shuts down the Kafka producer
func (p *KafkaProducer) Close() {
	p.Producer.Close()
}
