package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func rabbitMQDial(url string) (*amqp.Connection, error) {
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn, err
		}
		time.Sleep(time.Second * 3)
	}

}

func connectToQueue(conn *amqp.Connection) NotificationSource {
	if conn == nil {
		mqconn, _ = rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
		log.Println("Failed to connect to queue")
	}
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

	return NotificationSource{q: q, conn: conn, ch: ch}
}

type NotificationSource struct {
	q    amqp.Queue
	conn *amqp.Connection
	ch   *amqp.Channel
}

type NotificationSender struct {
	src <-chan amqp.Delivery
}

func (source NotificationSender) Get() string {
	var e string
	json.Unmarshal((<-source.src).Body, &e)
	return e
}
func (source *NotificationSource) SendUserNotification(toUser Username, notification string) error {
	s, _ := json.Marshal(notification)
	ch, err := source.conn.Channel()
	if err != nil {
		return err
	}
	return ch.Publish(
		"",                         // Exchange name
		"users"+"#"+string(toUser), // Queue name
		false,                      // Mandatory
		false,                      // Immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(s),
		},
	)

}

func (source *NotificationSource) postNotification(d Doc, notification string) error {
	s, _ := json.Marshal(notification)
	ch, err := source.conn.Channel()
	if err != nil {
		return err
	}
	return ch.Publish(
		"",                            // Exchange name
		d.CollectionName()+"#"+d.Id(), // Queue name
		false,                         // Mandatory
		false,                         // Immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(s),
		},
	)

}

func trackNotifications(d Doc, c *amqp.Connection) NotificationSender {
	connectToQueue(mqconn)
	if c == nil {
		log.Println(c, "connection was nil")
	}
	ch, err := c.Channel()
	if err != nil {
		log.Println("failed to create the channel")
	}
	// Just in case no notifications are sent yet
	_, err = ch.QueueDeclare(
		d.CollectionName()+"#"+d.Id(), // Queue name
		false,                         // Durable
		false,                         // Delete when unused
		false,                         // Exclusive
		false,                         // No-wait
		nil,                           // Arguments
	)

	failOnError(err, "Failed to open a channel")
	source, err := ch.Consume(d.CollectionName()+"#"+d.Id(), "", false, false, false, false, nil)
	failOnError(err, "Failed to open a channel")
	return NotificationSender{src: source}
}
