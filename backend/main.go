package main

import (
	"encoding/json"
	"log"
	"net/http"

	"context"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	dbHOST    = "10.9.0.3"
	mqHOST    = "10.9.0.10"
	redisHOST = "10.9.0.9"
)

// const (
// 	dbHOST    = "localhost"
// 	mqHOST    = "localhost"
// 	redisHOST = "localhost"
// )

var (
	db *dataStore
	// db       *mongo.Client
	// ctx      context.Context
	mycache  *redis.Client
	cacheCTX context.Context
	mqconn   *amqp.Connection
)

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
	MarshalBinary() ([]byte, error)
	isDoc() bool
	CollectionName() string
	Id() string
}

func Cache[D Doc](doc D, mycache *redis.Client, ctx context.Context) error {
	if mycache == nil {
		mycache, ctx = setupCache()
	}
	log.Println("cacheing value ", doc)
	if err := mycache.Set(cacheCTX, doc.CollectionName()+"#"+doc.Id(), doc, 0).Err(); err != nil {
		return err
	}
	return nil

}
func CheckCache[D Doc](d D, mycache *redis.Client, ctx context.Context) (D, error) {
	if mycache == nil {
		mycache, ctx = setupCache()
	}

	var wanted D
	got, err := mycache.Get(ctx, d.CollectionName()+"#"+d.Id()).Bytes()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return wanted, err
	} else if err == redis.Nil {
		log.Println("cache hit for", d)
	}
	log.Printf("we also got %s\n", string(got))
	err = json.Unmarshal(got, &wanted)
	if err != nil {
		log.Println(err)
	}
	return wanted, err
}

func main() {
	// ------- Connect to RabbitMQ server -------
	mqconn, _ = rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
	if mqconn == nil {
		log.Println("mq connection failed")
	}

	// ------- Setup Database --------
	db = connectToDB()
	defer db.Disconnect()

	// ------- Setup Cache ----
	mycache, cacheCTX = setupCache()

	// -------- Setup endpoints -----
	log.Println("Backend Is Now Running")
	http.HandleFunc("/", Signup)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/delete", DeleteAccount)
	http.HandleFunc("/post", NewPost)
	http.HandleFunc("/newfriend", NewFriend)
	http.HandleFunc("/stats", AccountStats)
	http.HandleFunc("/comment", NewComment)
	http.HandleFunc("/notifications", ReadNotifications)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
