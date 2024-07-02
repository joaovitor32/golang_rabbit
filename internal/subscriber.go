package rabbit_mq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQSubscriber struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	msgs       <-chan amqp.Delivery
	ready      chan<- struct{}
	done       chan struct{}
}

func NewRabbitMQSubscriber(amqpURI string, ready chan <-struct{}) (*RabbitMQSubscriber, error) {
	subscriber := &RabbitMQSubscriber{ready: ready}

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}
	subscriber.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	subscriber.channel = ch

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	subscriber.queue = q

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	// Start a goroutine to consume messages continuously
	go subscriber.consumeMessages(msgs)

	return subscriber, nil
}

func (s *RabbitMQSubscriber) consumeMessages(msgs <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Channel closed, stopping message consumption")
				return
			}
			fmt.Printf("[CONSUMER] Received a message: %s\n", msg.Body)

			// Process the message here

		case <-s.done:
			log.Println("Stopping message consumption")
			return
		}
	}
}

func (s *RabbitMQSubscriber) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.connection != nil {
		s.connection.Close()
	}
}

func (s *RabbitMQSubscriber) WaitForReady() {
	s.ready <- struct{}{}
}
