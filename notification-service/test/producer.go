package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Configure the producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Create the producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create a test event
	event := Event{
		ID:        uuid.New().String(),
		Type:      "user_registration",
		UserID:    "user123",
		Data:      map[string]string{"email": "test@example.com"},
		Timestamp: time.Now(),
	}

	// Marshal the event
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("Failed to marshal event: %v", err)
	}

	// Create the message
	msg := &sarama.ProducerMessage{
		Topic: "user-events",
		Value: sarama.StringEncoder(eventBytes),
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent successfully! Partition: %d, Offset: %d\n", partition, offset)
}
