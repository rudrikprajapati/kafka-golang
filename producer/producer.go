package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	kafkaTopic = "report-tasks"
)

var kafkaBroker = os.Getenv("KAFKA_BROKER")

type UploadData struct {
	UserID string `json:"user_id" binding:"required"`
	Data   string `json:"data" binding:"required"`
}

func main() {
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer producer.Close()

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send message: %v", err)
		}
	}()

	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		var upload UploadData
		if err := c.ShouldBindJSON(&upload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		start := time.Now()
		msg := &sarama.ProducerMessage{
			Topic: kafkaTopic,
			Key:   sarama.StringEncoder(upload.UserID),
			Value: sarama.StringEncoder(upload.Data),
		}
		producer.Input() <- msg

		select {
		case <-producer.Successes():
			log.Printf("Sent message for user %s in %v", upload.UserID, time.Since(start))
		case <-time.After(100 * time.Millisecond):
			log.Printf("Timeout waiting for send confirmation for user %s", upload.UserID)
		}

		c.JSON(http.StatusOK, gin.H{"status": "Processing..."})
	})

	log.Println("Producer starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
