package main

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

var (
	kafkaTopic  = "test-topic"
	kafkaBroker = os.Getenv("KAFKA_BROKER")
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func main() {
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092" // Default
	}

	consumer, err := NewKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer consumer.consumer.Close()

	consumeMessages(consumer)
}

func NewKafkaConsumer() (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{kafkaBroker}, config)
	if err != nil {
		log.Printf("Failed to connect to Kafka... (%v)", err)
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer}, nil
}

func consumeMessages(kc *KafkaConsumer) {
	partitionConsumer, err := kc.consumer.ConsumePartition(kafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Failed to start consumer: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message: %s", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)
		}
	}
}
