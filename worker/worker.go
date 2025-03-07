package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    "github.com/Shopify/sarama"    // Kafka client library
)

// connectConsumer creates and configures a Kafka consumer with retries
func connectConsumer(brokers []string) (sarama.Consumer, error) {
    config := sarama.NewConfig()                    // Create Kafka config
    config.Consumer.Return.Errors = true            // Enable error reporting
    config.Consumer.Retry.Backoff = time.Second * 2 // Retry delay
    config.Consumer.Retry.Max = 5                   // Max retry attempts
    config.Net.DialTimeout = time.Second * 10       // Connection timeout
    
    // Retry connection with backoff
    var consumer sarama.Consumer
    var err error
    for retries := 0; retries < 3; retries++ {
        consumer, err = sarama.NewConsumer(brokers, config)
        if err == nil {
            return consumer, nil
        }
        log.Printf("Failed to connect to Kafka, attempt %d/3: %v\n", retries+1, err)
        time.Sleep(time.Second * 2)
    }
    return nil, fmt.Errorf("failed to create consumer after retries: %w", err)
}

func main() {
    // Kafka configuration
    brokers := []string{"localhost:9092"}                   // Updated Kafka broker address
    topic := "comments"

    log.Println("Connecting to Kafka...")
    
    // Create consumer with retries
    worker, err := connectConsumer(brokers)
    if err != nil {
        log.Fatalf("Consumer creation failed: %v", err)
    }
    defer worker.Close()

    // Create partition consumer
    consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        log.Fatalf("Partition consumer creation failed: %v", err)
    }
    defer consumer.Close()

    log.Printf("Consumer started successfully on topic: %s\n", topic)
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    msgCount := 0
    doneCh := make(chan struct{})

    // Message consumption loop
    go func() {
        for {
            select {
            case err := <-consumer.Errors():
                log.Printf("Consumer error: %v\n", err)
                
            case msg := <-consumer.Messages():
                msgCount++
                log.Printf("Received message %d:\n  Topic: %s\n  Partition: %d\n  Offset: %d\n  Value: %s\n",
                    msgCount, msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
                
            case sig := <-sigChan:
                log.Printf("Caught signal %v: terminating\n", sig)
                doneCh <- struct{}{}
                return
            }
        }
    }()

    <-doneCh
    log.Println("Gracefully shutting down consumer...")
}