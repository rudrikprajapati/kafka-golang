package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

var (
	kafkaTopic  = "test-topic"
	kafkaBroker = os.Getenv("KAFKA_BROKER")
)

type Message struct {
	Content string `json:"content" binding:"required"`
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func main() {
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092" // Default
	}

	producer, err := NewKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer after retries: %v", err)
	}
	defer producer.producer.Close()

	r := gin.Default()
	r.POST("/produce", func(c *gin.Context) {
		var msg Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = producer.SendMessage(msg.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "message sent"})
	})
	r.Run(":8080")
}

func NewKafkaProducer() (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) SendMessage(content string) error {
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(content),
	}
	_, _, err := kp.producer.SendMessage(msg)
	return err
}
