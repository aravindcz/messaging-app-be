package main

import (
        "database/sql"
        "encoding/json"
        "log"
        "net/http"
        "os"
        "strconv"

        "github.com/gin-gonic/gin"
        _ "github.com/lib/pq"
        "github.com/streadway/amqp"
)

// Database and RabbitMQ connection details
var (
        db  *sql.DB
        rmq *amqp.Connection
)

// Message struct
type Message struct {
        ID         int    `json:"id"`
        SenderID   int `json:"sender_id" binding:"required"`
        ReceiverID int `json:"receiver_id" binding:"required"`
        Content    string `json:"content" binding:"required"`
        Timestamp  string `json:"timestamp"`
        Read       bool   `json:"read"`
}

func main() {
        // Initialize database connection
        var err error
        db, err = sql.Open("postgres", "postgres://postgres:password@postgres:5432/messaging?sslmode=disable")
        if err != nil {
                log.Fatal("Error connecting to the database: ", err)
        }
        defer db.Close()

        // Connect to RabbitMQ
        rmq, err = amqp.Dial("amqp://admin:adminpassword@rabbitmq:5672/")
        if err != nil {
                log.Fatal("Failed to connect to RabbitMQ: ", err)
        }
        defer rmq.Close()

        // Setup router
        router := gin.Default()
        router.POST("/messages", sendMessage)
        router.GET("/messages", getMessages)
        router.PATCH("/messages/:message_id/read", markMessageAsRead)

        // Start API
        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }
        router.Run(":" + port)
}

// Send Message (Push to Queue)
func sendMessage(c *gin.Context) {
        var msg Message
        if err := c.ShouldBindJSON(&msg); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
                return
        }

        ch, err := rmq.Channel()
        if err != nil {
                log.Println("Failed to open channel:", err)
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
                return
        }
        defer ch.Close()

        // Declare queue
        q, err := ch.QueueDeclare("message_queue", true, false, false, false, nil)
        if err != nil {
                log.Println("Queue declaration failed:", err)
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
                return
        }

        // Publish message to queue
        body, _ := json.Marshal(msg)
        err = ch.Publish("", q.Name, false, false, amqp.Publishing{
                ContentType: "application/json",
                Body:        body,
        })
        if err != nil {
                log.Println("Failed to publish message:", err)
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send message"})
                return
        }

        c.JSON(http.StatusAccepted, gin.H{"message": "Message sent to queue"})
}

// Get Conversation History
func getMessages(c *gin.Context) {
	// Check if both user1 and user2 parameters exist in the request
	user1, user1Exists := c.GetQuery("user1")
	user2, user2Exists := c.GetQuery("user2")

	if !user1Exists || !user2Exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Both user1 and user2 query parameters are required"})
		return
	}

	// Ensure values are not empty
	if user1 == "" || user2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Both user1 and user2 must have valid values"})
		return
	}

        rows, err := db.Query(`
                SELECT id, sender_id, receiver_id, content, timestamp, read 
                FROM messages 
                WHERE (sender_id = $1 AND receiver_id = $2) 
                   OR (sender_id = $2 AND receiver_id = $1) 
                ORDER BY timestamp ASC`, user1, user2)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
                return
        }
        defer rows.Close()

        var messages []Message
        for rows.Next() {
                var msg Message
                err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read)
                if err != nil {
                        log.Println("Error scanning row:", err)
                        continue
                }
                messages = append(messages, msg)
        }

        c.JSON(http.StatusOK, messages)
}

// Mark Message as Read
func markMessageAsRead(c *gin.Context) {
        messageID := c.Param("message_id")

        id, Interr := strconv.Atoi(messageID)
        if Interr != nil {
                c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid message_id"})
                return
        }

        _, err := db.Exec("UPDATE messages SET read = true WHERE id = $1", messageID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"status": "read"})
}