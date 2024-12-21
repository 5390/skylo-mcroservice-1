package queue

import (
	"context"
	"log"

	"microservice-1/config"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(config config.QueueConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Broker},
		Topic:   config.Topic,
		GroupID: config.GroupID,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) Messages() <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				continue
			}
			out <- string(msg.Value)
		}
	}()
	return out
}
