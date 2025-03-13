package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"notification-service/config"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Event represents a user event that will be processed
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// Notification represents a notification to be sent
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Channel   string    `json:"channel"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	cfg           *config.Config
	consumerGroup sarama.ConsumerGroup
	producer      sarama.SyncProducer
)

func setupKafka() error {
	// Configure Kafka consumer group
	kafkaConfig := sarama.NewConfig()

	// Set consumer group configuration
	switch cfg.Kafka.Consumer.Strategy {
	case "roundrobin":
		kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	}

	switch cfg.Kafka.Consumer.InitialOffset {
	case "newest":
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	case "oldest":
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// Create consumer group
	var err error
	consumerGroup, err = sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, kafkaConfig)
	if err != nil {
		return fmt.Errorf("error creating consumer group: %v", err)
	}

	// Configure and create producer
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = cfg.Kafka.Producer.ReturnSuccesses
	producerConfig.Producer.Return.Errors = cfg.Kafka.Producer.ReturnErrors
	producer, err = sarama.NewSyncProducer(cfg.Kafka.Brokers, producerConfig)
	if err != nil {
		return fmt.Errorf("error creating producer: %v", err)
	}

	return nil
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			var event Event
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			// Process the event and create notifications
			processEvent(event)

			// Mark message as processed
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func processEvent(event Event) {
	var notifications []Notification

	switch event.Type {
	case "user_registration":
		notifications = []Notification{
			{
				ID:        uuid.New().String(),
				UserID:    event.UserID,
				Type:      "welcome",
				Message:   "Welcome to our platform!",
				Channel:   "email",
				Status:    "pending",
				Timestamp: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				UserID:    event.UserID,
				Type:      "profile_update",
				Message:   "Your profile has been updated successfully",
				Channel:   "push",
				Status:    "pending",
				Timestamp: time.Now(),
			},
		}
	case "user_login":
		notifications = []Notification{
			{
				ID:        uuid.New().String(),
				UserID:    event.UserID,
				Type:      "login",
				Message:   "Hello, welcome back!",
				Channel:   "push",
				Status:    "pending",
				Timestamp: time.Now(),
			},
		}
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return
	}

	// Send notifications to Kafka
	for _, notification := range notifications {
		notificationBytes, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Error marshaling notification: %v", err)
			continue
		}

		msg := &sarama.ProducerMessage{
			Topic: cfg.Kafka.Topics.Notifications,
			Value: sarama.StringEncoder(notificationBytes),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Error sending notification: %v", err)
			continue
		}

		log.Printf("Notification sent - Partition: %d, Offset: %d", partition, offset)
	}
}

func main() {
	// Load configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup Kafka
	if err := setupKafka(); err != nil {
		log.Fatalf("Failed to setup Kafka: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		}
		if err := producer.Close(); err != nil {
			log.Printf("Error closing producer: %v", err)
		}
	}()

	// Setup Gin router
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start HTTP server
	go func() {
		if err := router.Run(cfg.HTTP.Port); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	consumer := Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{cfg.Kafka.Topics.UserEvents}, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.ready
	log.Println("Sarama consumer up and running!...")

	// Wait for interrupt signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Println("Initiating shutdown...")
	cancel()
	wg.Wait()
}
