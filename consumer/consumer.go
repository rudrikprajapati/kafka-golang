package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

const (
	kafkaTopic = "report-tasks"
	groupID    = "report-workers"
)

var kafkaBroker = os.Getenv("KAFKA_BROKER")

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		data := string(msg.Value)
		log.Printf("Starting report for user %s with data: %s", userID, data)
		start := time.Now()

		time.Sleep(10 * time.Second)

		log.Printf("Finished report for user %s in %v", userID, time.Since(start))
		session.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup([]string{kafkaBroker}, groupID, config)
	if err != nil {
		log.Fatalf("Failed to start consumer group: %v", err)
	}
	defer client.Close()

	handler := ConsumerGroupHandler{}
	ctx := context.Background()

	log.Println("Consumer starting...")
	for {
		err := client.Consume(ctx, []string{kafkaTopic}, handler)
		if err != nil {
			log.Printf("Error consuming: %v, retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}
	}
}
