package config

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitRabbitMQ(url string) error {
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func Publish(queue string, body []byte) error {
	_, err := ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}

// Consume messages from RabbitMQ
func Consume(queue string) {
	msgs, err := ch.Consume(
		queue, // queue name
		"",    // consumer tag
		true,  // auto-acknowledge
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
		// Process message here, such as saving to DB
	}
}
