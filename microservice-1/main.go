package main

import (
	"log"
	"microservice-1/config"
	"microservice-1/queue"
	"microservice-1/retry"
)

func main() {
	// Initialize configurations
	cfg := config.LoadConfig()

	// Start consuming messages from the queue
	log.Println("Starting Microservice-1...")
	consumer := queue.NewConsumer(cfg.QueueConfig)
	retryHandler := retry.NewRetryHandler(cfg.RetryConfig)

	for message := range consumer.Messages() {
		go func(msg string) {
			err := retryHandler.ProcessMessage(msg)
			if err != nil {
				log.Printf("Failed to process message: %s, Error: %v\n", msg, err)
			}
		}(message)
	}
}
