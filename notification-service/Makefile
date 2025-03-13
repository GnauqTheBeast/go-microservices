.PHONY: run build test clean docker-up docker-down

# Default target
all: build

# Build the application
build:
	go build -o bin/notification-service main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Kafka and Zookeeper
docker-up:
	docker-compose up -d

# Stop Kafka and Zookeeper
docker-down:
	docker-compose down

# Install dependencies
deps:
	go mod tidy

# Run all setup steps
setup: deps docker-up 