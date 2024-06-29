package main

import (
	"fmt"
	"sync"
	"time"

	rabbit_mq "github.com/joaovitor32/golang_rabbit/internal"
)

func main() {
	var wg sync.WaitGroup
	started := make(chan struct{})

	wg.Add(2)

	// Subscribe goroutine
	go func() {
		defer wg.Done()
		fmt.Println("Starting Subscribe...")
		rabbit_mq.Subscribe(started)
		fmt.Println("Subscribe started")
	}()

	// Publish goroutine
	go func() {
		defer wg.Done()
		// Wait for the signal that Subscribe has started
		<-started
		fmt.Println("Starting Publish after Subscribe...")
		// Optional delay to ensure Subscribe is fully ready
		time.Sleep(1 * time.Second)
		rabbit_mq.Publish()
		fmt.Println("Publish started")
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	fmt.Println("Both Publish and Subscribe have finished.")
}