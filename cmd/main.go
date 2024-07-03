package main

import (
	"fmt"
	"log"
	"sync"

	rabbit_mq "github.com/joaovitor32/golang_rabbit/internal"
)

const URI = "amqp://guest:guest@localhost:5672/"
var messages = []string{"Never", "gonna", "give", "you", "up", "never", "gonna", "let", "you", "down", "and", "hurt", "you"}

func main() {
	var wg sync.WaitGroup
	ready := make(chan struct{}, 1) // Buffered channel

	wg.Add(2)

	// Subscribe goroutine
	go func() {
		defer wg.Done()
		fmt.Println("Starting Subscribe...")

		subscriber, err := rabbit_mq.NewRabbitMQSubscriber(URI, ready)
		if err != nil {
			log.Fatalf("Failed to create RabbitMQ subscriber: %v", err)
		}
		// Signal that subscriber has started

		subscriber.WaitForReady()

		defer subscriber.Close()
	}()

	// Publish goroutine
	go func() {
		defer wg.Done()

		// Wait for the signal that Subscribe has started
		fmt.Println("Starting Publish after Subscribe...")
		ready <- struct{}{}
		
		publisher, err := rabbit_mq.NewRabbitMQPublisher(URI)
		if err != nil {
			log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
		}

		defer publisher.Close()

		fmt.Println("Publish started")
		err = publisher.Publish(messages)
		if err != nil {
			log.Fatalf("Failed to publish messages: %v", err)
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	fmt.Println("Both Publish and Subscribe have finished.")
}
