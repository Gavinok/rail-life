package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func rabbitMQDial(url string) (*amqp.Connection, error) {
	for {
		conn, err := amqp.Dial("amqp://guest:guest@10.9.0.10:5672/")
		if err == nil {
			return conn, err
		}
		time.Sleep(time.Second * 3)
	}

}
func connectToQueue(conn *amqp.Connection) (amqp.Queue, *amqp.Channel) {
	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Declare a queue
	q, err := ch.QueueDeclare(
		"Notifications", // Queue name
		false,           // Durable
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q, ch
}

func postNotification(user Username, notification string, q *amqp.Channel) error {
	s, _ := json.Marshal(notification)
	return q.Publish(
		"",                            // Exchange name
		"Notifications#"+string(user), // Queue name
		false,                         // Mandatory
		false,                         // Immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(s),
		},
	)

}

func trackNotifications(user Username, c *amqp.Connection) <-chan amqp.Delivery {
	ch, err := c.Channel()

	// Just in case no notifications are sent yet

	_, err = ch.QueueDeclare(
		"Notifications#"+string(user), // Queue name
		false,                         // Durable
		false,                         // Delete when unused
		false,                         // Exclusive
		false,                         // No-wait
		nil,                           // Arguments
	)

	failOnError(err, "Failed to open a channel")
	d, err := ch.Consume("Notifications#"+string(user), "", false, false, false, false, nil)
	failOnError(err, "Failed to open a channel")
	return d
}
