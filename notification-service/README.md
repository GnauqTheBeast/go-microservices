# Kafka Power Demonstration - Notification Service

This service demonstrates the power of Apache Kafka in a microservices architecture by implementing a real-time notification system. It showcases several key Kafka features and best practices.

## Key Features Demonstrated

1. **Event-Driven Architecture**
   - Consumes user events from Kafka
   - Processes events in real-time
   - Generates notifications based on event types

2. **Kafka Consumer Groups**
   - Implements consumer group for load balancing
   - Handles partition rebalancing
   - Ensures message processing reliability

3. **Message Production**
   - Produces notifications to Kafka topics
   - Implements synchronous message production
   - Handles message delivery guarantees

4. **Fault Tolerance**
   - Graceful shutdown handling
   - Error handling and recovery
   - Message processing retries

5. **Configuration Management**
   - Uses Viper for flexible configuration
   - Supports multiple configuration sources
   - Environment variable overrides
   - Default values

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## Configuration

The service uses Viper for configuration management. Configuration can be provided through:

1. YAML configuration file (`config/config.yaml`)
2. Environment variables (prefixed with `APP_`)
3. Command-line flags

### Configuration Options

```yaml
kafka:
  brokers:
    - localhost:9092
  consumerGroup: notification-group
  topics:
    userEvents: user-events
    notifications: notifications
  consumer:
    initialOffset: oldest  # or newest
    strategy: roundrobin   # or range
  producer:
    returnSuccesses: true
    returnErrors: true

http:
  port: :8080
```

### Environment Variables

All configuration options can be overridden using environment variables:

```bash
APP_KAFKA_BROKERS=["localhost:9092"]
APP_KAFKA_CONSUMERGROUP=notification-group
APP_KAFKA_TOPICS_USEREVENTS=user-events
APP_KAFKA_TOPICS_NOTIFICATIONS=notifications
APP_KAFKA_CONSUMER_INITIALOFFSET=oldest
APP_KAFKA_CONSUMER_STRATEGY=roundrobin
APP_KAFKA_PRODUCER_RETURNSUCCESSES=true
APP_KAFKA_PRODUCER_RETURNERRORS=true
APP_HTTP_PORT=:8080
```

## Setup and Running

1. Start Kafka and Zookeeper:
```bash
docker-compose up -d
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the service:
```bash
go run main.go
```

## Testing the Service

1. The service listens on port 8080 for health checks
2. It consumes messages from the `user-events` topic
3. It produces notifications to the `notifications` topic

## Architecture Overview

```
[User Service] -> [Kafka: user-events] -> [Notification Service] -> [Kafka: notifications]
```

The service demonstrates:
- Event streaming
- Real-time processing
- Message transformation
- Topic-to-topic communication
- Configuration management

## Environment Variables

- `KAFKA_BROKER`: Kafka broker address (default: localhost:9092)

## Best Practices Implemented

1. **Consumer Group Management**
   - Proper setup and cleanup
   - Partition rebalancing strategy
   - Message acknowledgment

2. **Producer Configuration**
   - Synchronous message production
   - Error handling
   - Message delivery guarantees

3. **Error Handling**
   - Graceful error recovery
   - Logging and monitoring
   - Circuit breaking

4. **Resource Management**
   - Proper connection cleanup
   - Context cancellation
   - Graceful shutdown

5. **Configuration Management**
   - Centralized configuration
   - Multiple configuration sources
   - Environment variable support
   - Default values 