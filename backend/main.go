package main

import (
	"encoding/json"
	"log"
	"net/http"

	"context"
	"time"

	// "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// const dbHOST = "10.9.0.3"
// const mqHOST = "10.9.0.10"
const dbHOST = "localhost"
const mqHOST = "localhost"
const redisHOST = "localhost"

var db *mongo.Client
var ctx context.Context

var mycache *redis.Client
var cacheCTX context.Context

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setupCache() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		DB:       0,
		Password: "",
		Addr:     redisHOST + ":6379",
	})

	cacheCTX := context.Background()
	return client, cacheCTX
}

type Doc interface {
	isDoc() bool
	CollectionName() string
	Id() string
}

func Cache[D Doc](doc D, mycache *redis.Client, ctx context.Context) error {
	// data, err := json.Marshal(doc)
	// if err != nil {
	// 	return err
	// }
	log.Println("cacheing value ", doc)
	if err := mycache.Set(cacheCTX, doc.CollectionName()+"#"+doc.Id(), doc, 0).Err(); err != nil {
		return err
	}
	return nil

}
func CheckCache[D Doc](d D, mycache *redis.Client, ctx context.Context) (D, error) {
	var wanted D
	got, err := mycache.Get(ctx, d.CollectionName()+"#"+d.Id()).Bytes()
	if err != nil {
		log.Println(err)
	}
	log.Printf("we also got %s\n", string(got))
	err = json.Unmarshal(got, &wanted)
	if err != nil {
		log.Println(err)
	}
	return wanted, err
}

// TODO for each friend look up any remaining posts
func main() {
	// ------- Connect to RabbitMQ server -------
	conn, err := rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
	defer conn.Close()
	_, ch := connectToQueue(conn)
	defer ch.Close()

	// ------- Setup Database --------
	db, ctx = connectToDB()
	defer db.Disconnect(ctx)

	// ------- Setup Cache ----
	mycache, cacheCTX = setupCache()
	// TODO  Now we need to track this users notifications

	failOnError(err, "failed to read user")
	// -------- Setup endpoints -----
	u, _ := create_user("gavin", Username("gavinok"), "tmp", "", db, ctx)
	time.Sleep(2 * time.Second)
	err = Cache(u, mycache, cacheCTX)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(2 * time.Second)
	u2, err := CheckCache(u, mycache, cacheCTX)
	if err != nil {
		log.Println(err)
	}
	log.Println("we got ", *u2)
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

	// newNotChan, _ := conn.Channel()
	// Consume messages from the queue
	// for {
	// 	e := postNotification(u2.Username, "hello", newNotChan)
	// 	if e != nil {
	// 		print(e)
	// 	}
	// 	time.Sleep(10 * time.Millisecond)
	// 	m := <-msgs
	// 	var s string
	// 	json.Unmarshal(m.Body, &s)
	// 	print(s)

	// }
}
