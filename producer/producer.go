package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/Shopify/sarama"       // Kafka client library
    "github.com/gofiber/fiber/v2"     // Web framework
)

type Comment struct {
    Text string `form:"text" json:"text" binding:"required"`    // Text field with validation tag
}

func main() {
    app := fiber.New()                         // Initialize web server
    api := app.Group("/api/v1")               // Create versioned API group
    api.Post("/comments", createComment)       // Register comment endpoint
    log.Fatal(app.Listen(":3000"))            // Start server on port 3000
}

func createProducer(brokersUrl []string) (sarama.SyncProducer, error) {
    config := sarama.NewConfig()                              // Create Kafka config
    
    config.Producer.RequiredAcks = sarama.WaitForAll         // Strongest consistency
    config.Producer.Retry.Max = 5                            // Retry attempts
    config.Producer.Return.Successes = true                  // Required for sync producer
    
    conn, err := sarama.NewSyncProducer(brokersUrl, config) // Create producer
    if err != nil {
        return nil, fmt.Errorf("failed to create producer: %w", err)
    }
    return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
    brokersUrl := []string{"localhost:29092"}                // Kafka broker address
    
    producer, err := createProducer(brokersUrl)             // Initialize producer
    if err != nil {
        return fmt.Errorf("failed to initialize producer: %w", err)
    }
    defer producer.Close()                                   // Ensure cleanup
    
    msg := &sarama.ProducerMessage{                         // Create message
        Topic: topic,
        Value: sarama.StringEncoder(message),
    }
    
    partition, offset, err := producer.SendMessage(msg)      // Send to Kafka
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    
    log.Printf("Message stored in topic(%s)/partition(%d)/offset(%d)\n",  // Log success
        topic, partition, offset)
    return nil
}

func createComment(c *fiber.Ctx) error {
    cmt := new(Comment)                                      // Initialize comment
    
    if err := c.BodyParser(cmt); err != nil {               // Parse request body
        log.Printf("Failed to parse request body: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request body",
            "error":   err.Error(),
        })
    }
    
    if cmt.Text == "" {                                     // Validate input
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Comment text cannot be empty",
        })
    }
    
    cmtInBytes, err := json.Marshal(cmt)                    // Convert to JSON
    if err != nil {
        log.Printf("Failed to marshal comment: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to process comment",
        })
    }
    
    if err := PushCommentToQueue("comments", cmtInBytes); err != nil {  // Send to Kafka
        log.Printf("Failed to push comment to queue: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to queue comment",
        })
    }

    return c.JSON(fiber.Map{                                // Return success
        "success": true,
        "message": "Comment pushed successfully",
        "comment": cmt,
    })
}