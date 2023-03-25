package main

import (
	"encoding/json"
	"log"
	"net/http"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// const dbHOST = "10.9.0.3"
// const mqHOST = "10.9.0.10"
const dbHOST = "localhost"
const mqHOST = "localhost"

var db *mongo.Client
var ctx context.Context

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// TODO for each friend look up any remaining posts
func main() {
	// Connect to RabbitMQ server
	conn, err := rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
	defer conn.Close()
	_, ch := connectToQueue(conn)
	defer ch.Close()

	user := Username("gavin")

	// Setup Database
	db, ctx = connectToDB()
	defer db.Disconnect(ctx)

	// User Example
	u, err := create_user("gavin", user, "gavinfreeborn@gmail.com", "1997", db, ctx)
	failOnError(err, "failed to create user")
	u2, err := read_user(u.Username, db, ctx)
	log.Println(u2)
	// TODO  Now we need to track this users notifications

	// Now we need to track this users notifications
	msgs := trackNotifications(u2.Username, conn)

	failOnError(err, "failed to read user")
	// TODO create comment on a post
	// TODO like a post
	// TODO like comments
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/delete", DeleteAccount)
	http.HandleFunc("/post", NewPost)
	http.HandleFunc("/stats", AccountStats)
	// TODO incomplete since I don't have a way to look up posts
	// with this now
	http.HandleFunc("/comment", NewComment)
	log.Fatal(http.ListenAndServe(":8000", nil))

	newNotChan, _ := conn.Channel()
	// Consume messages from the queue
	for {
		e := postNotification(u2.Username, "hello", newNotChan)
		if e != nil {
			print(e)
		}
		time.Sleep(10 * time.Millisecond)
		m := <-msgs
		var s string
		json.Unmarshal(m.Body, &s)
		print(s)

	}
}
