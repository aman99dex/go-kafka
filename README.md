# Go-Kafka Comment Service

This project demonstrates how to implement a comment service using Go and Apache Kafka, featuring a REST API built with the Fiber web framework.

## What is Kafka?

Apache Kafka is a distributed streaming platform that is used to build real-time data pipelines and streaming applications. It is designed to handle large volumes of data with high throughput and low latency. Kafka is used for various use cases such as messaging, website activity tracking, log aggregation, stream processing, and more.

## Features of Kafka

- **Scalability**: Kafka can scale horizontally by adding more brokers to the cluster.
- **Durability**: Kafka persists messages on disk and replicates them across multiple brokers to ensure data durability.
- **High Throughput**: Kafka can handle high throughput of messages with low latency.
- **Fault Tolerance**: Kafka is designed to be fault-tolerant by replicating data across multiple brokers.

## Architecture Overview

The project consists of two main components:
1. **Producer Service**: A REST API that accepts comments and publishes them to Kafka
2. **Consumer Service**: A service that reads comments from Kafka and processes them

## Technical Stack

- **Go**: Programming language (v1.22.6+)
- **Fiber**: Web framework for REST API
- **Apache Kafka**: Message broker for async processing
- **Docker**: Containerization platform

## API Endpoints

### Comments API

```
POST /api/v1/comments
Content-Type: application/json

{
    "text": "Your comment text here"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Comment pushed successfully",
    "comment": {
        "text": "Your comment text here"
    }
}
```

## Project Structure

```
/go-kafka
│
├── producer/
│   ├── producer.go    # REST API implementation
│   └── kafka.go       # Kafka producer implementation
│
├── consumer/
│   └── main.go        # Kafka consumer implementation
│
├── docker-compose.yml # Docker compose for local development
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22.6 or higher
- Kafka 2.8.0 or higher
- Docker and Docker Compose

### Local Development Setup

1. Start Kafka using Docker Compose:
```sh
docker-compose up -d
```

2. Install dependencies:
```sh
go mod tidy
```

3. Start the producer service:
```sh
cd producer
go run *.go
```

4. Start the consumer service:
```sh
cd consumer
go run main.go
```

### Testing the API

You can test the API using curl:

```sh
curl -X POST http://localhost:3000/api/v1/comments \
  -H "Content-Type: application/json" \
  -d '{"text":"Hello, World!"}'
```

## Error Handling

The API implements proper error handling:
- 400 Bad Request: Invalid request body
- 500 Internal Server Error: Server-side processing errors

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
