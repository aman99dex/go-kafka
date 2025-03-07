package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/Shopify/sarama"    // Kafka client library
)

// connectConsumer creates and configures a Kafka consumer
func connectConsumer(brokers []string) (sarama.Consumer, error) {
    config := sarama.NewConfig()                    // Create Kafka config
    config.Consumer.Return.Errors = true            // Enable error reporting
    
    consumer, err := sarama.NewConsumer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer: %w", err)
    }
    return consumer, nil
}

func main() {
    topic := "comments"                            // Kafka topic to consume from
    worker, err := connectConsumer([]string{"localhost:29092"})
    if err != nil {
        log.Fatal("Consumer creation failed:", err)
    }
    defer worker.Close()                           // Ensure cleanup on exit

    consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        log.Fatal("Partition consumer creation failed:", err)
    }
    defer consumer.Close()                         // Ensure cleanup on exit

    log.Println("Consumer started successfully")    // Startup notification
    
    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    msgCount := 0                                  // Track processed messages
    doneCh := make(chan struct{})                 // Channel for shutdown signal

    // Start message consumption loop
    go func() {
        for {
            select {
            case err := <-consumer.Errors():
                log.Printf("Error: %v\n", err)     // Log consumer errors
                
            case msg := <-consumer.Messages():
                msgCount++
                log.Printf("Message %d | Topic: %s | Value: %s\n",
                    msgCount, msg.Topic, string(msg.Value))
                
            case <-sigChan:
                log.Println("Shutdown signal received")
                doneCh <- struct{}{}
                return
            }
        }
    }()

    <-doneCh                                       // Wait for shutdown signal
    log.Println("Gracefully shutting down...")
}