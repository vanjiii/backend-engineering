package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func send() {

	// setup socket connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// abstraction over the API
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("failed to open a channel: %s", err)
	}
	defer ch.Close()

	// delcare a queue to send to
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("failed to declare a queue: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Panicf("failed to publish a message: %s", err)
	}

	log.Printf("=>  [x] Sent %s\n", body)
}
