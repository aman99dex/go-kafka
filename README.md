# Go-Kafka Project

This project demonstrates how to use Apache Kafka with the Go programming language.

## What is Kafka?

Apache Kafka is a distributed streaming platform that is used to build real-time data pipelines and streaming applications. It is designed to handle large volumes of data with high throughput and low latency. Kafka is used for various use cases such as messaging, website activity tracking, log aggregation, stream processing, and more.

## Features of Kafka

- **Scalability**: Kafka can scale horizontally by adding more brokers to the cluster.
- **Durability**: Kafka persists messages on disk and replicates them across multiple brokers to ensure data durability.
- **High Throughput**: Kafka can handle high throughput of messages with low latency.
- **Fault Tolerance**: Kafka is designed to be fault-tolerant by replicating data across multiple brokers.

## Project Structure

```
/go-kafka
│
├── producer
│   └── main.go
│
├── consumer
│   └── main.go
│
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Kafka 2.8.0 or higher
- Docker (optional, for running Kafka locally)

### Running Kafka Locally

You can run Kafka locally using Docker. Here is a simple `docker-compose.yml` file to get you started:

```yaml
version: '2'
services:
    zookeeper:
        image: wurstmeister/zookeeper:3.4.6
        ports:
            - "2181:2181"
    kafka:
        image: wurstmeister/kafka:2.12-2.2.1
        ports:
            - "9092:9092"
        environment:
            KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
            KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

Run the following command to start Kafka:

```sh
docker-compose up -d
```

### Running the Producer

Navigate to the `producer` directory and run the following command:

```sh
go run main.go
```

### Running the Consumer

Navigate to the `consumer` directory and run the following command:

```sh
go run main.go
```

## Future Work

This README will be updated with more details as the project progresses. Stay tuned!
