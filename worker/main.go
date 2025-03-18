package main

import (
        "database/sql"
        "encoding/json"
        "log"

        _ "github.com/lib/pq"
        "github.com/streadway/amqp"
)

// Message struct
type Message struct {
        SenderID   string `json:"sender_id"`
        ReceiverID string `json:"receiver_id"`
        Content    string `json:"content"`
}

// Database connection string
const dbConnStr = "postgres://postgres:password@postgres:5432/messaging?sslmode=disable"

func main() {
        // Connect to RabbitMQ
        conn, err := amqp.Dial("amqp://admin:adminpassword@rabbitmq:5672/")
        if err != nil {
                log.Fatal("Failed to connect to RabbitMQ:", err)
        }
        defer conn.Close()

        ch, err := conn.Channel()
        if err != nil {
                log.Fatal("Failed to open a channel:", err)
        }
        defer ch.Close()

        // Connect to PostgreSQL
        db, err := sql.Open("postgres", dbConnStr)
        if err != nil {
                log.Fatal("Error connecting to the database:", err)
        }
        defer db.Close()

        // Declare queue
        q, err := ch.QueueDeclare("message_queue", true, false, false, false, nil)
        if err != nil {
                log.Fatal("Queue declaration failed:", err)
        }

        // Consume messages
        msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
        if err != nil {
                log.Fatal("Failed to consume messages:", err)
        }

        log.Println("Worker listening for messages...")
        for msg := range msgs {
                var message Message
                if err := json.Unmarshal(msg.Body, &message); err != nil {
                        log.Println("Failed to parse message:", err)
                        continue
                }

                // Insert into database
                _, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, content, timestamp, read) VALUES ($1, $2, $3, NOW(), false)",
                        message.SenderID, message.ReceiverID, message.Content)
                if err != nil {
                        log.Println("Failed to insert message:", err)
                }
        }
}