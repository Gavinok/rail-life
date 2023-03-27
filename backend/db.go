package main

import (
	"context"
	"errors"
	"log"

	"encoding/json"
	// "errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDB() (*mongo.Client, context.Context) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://admin:admin@" + dbHOST + ":27017")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		fmt.Println("Error connecting to DB")
		panic(err)
	}
	return client, ctx
}

type Username string
type user_doc struct {
	Name        string
	Username    Username
	Email       string
	DateOfBirth string
	Friends     []Username
	// list of post titles
	// Posts []string
	// list of comment contents
	// Comments []string
	// TODO maybe keep liked posts here
}

func (user_doc) isDoc() bool {
	return true
}
func (u user_doc) Id() string {
	return string(u.Username)
}
func (u user_doc) CollectionName() string {
	return "users"
}
func (i user_doc) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type post_doc struct {
	Title          string
	Content        string
	Author         Username
	DateOfCreation time.Time
	Likes          []Username
	Comments       []comment_doc
}

func (post_doc) isDoc() bool {
	return true
}
func (u post_doc) Id() string {
	return string(u.Title)
}
func (u post_doc) CollectionName() string {
	return "posts"
}
func (i post_doc) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type comment_doc struct {
	Content        string
	Author         Username
	DateOfCreation time.Time
	Likes          []Username
}

func (comment_doc) isDoc() bool {
	return true
}
func (u comment_doc) Id() string {
	return string(u.Content)
}
func (u comment_doc) CollectionName() string {
	return "comments"
}
func (i comment_doc) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

// TODO track user activity

// Different types of notifications
const (
	COMMENT_NOTIFICATION = iota
	LIKE_NOTIFICATION
	POST_NOTIFICATION
)

type NotificationType int8

type notification_doc struct {
	Type    NotificationType
	Content string
	Author  Username
}

const DATABASE = "social_media_app"

func read_user(user Username, db *mongo.Client, ctx context.Context) (*user_doc, error) {
	// Check the cache for efficinency reasons
	usr, err := CheckCache(user_doc{Username: user}, mycache, cacheCTX)
	if err == nil {
		log.Println("Cache hit")
		return &usr, nil
	} else {
		log.Println("error with cache ", err)
	}

	var result user_doc
	err = db.Database(DATABASE).Collection("users").FindOne(
		ctx, bson.D{{"username", string(user)}},
	).Decode(&result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
		return nil, errors.New("User does not yet exist")

	}

	Cache(result, mycache, cacheCTX)

	return &result, nil
}

func create_user(user string, username Username, email string, dateofbirth string, db *mongo.Client, ctx context.Context) (*user_doc, error) {

	u := user_doc{
		Name:        user,
		Username:    username,
		Email:       email,
		DateOfBirth: dateofbirth,
		Friends:     make([]Username, 0),
	}
	_, err := db.Database(DATABASE).Collection("users").InsertOne(context.TODO(), u)

	err = Cache(u, mycache, cacheCTX)
	CheckCache(u, mycache, cacheCTX)
	if err != nil {
		fmt.Println("Error adding user to db: ", err)
		panic(err)
	}

	return &u, nil

}
func (u *user_doc) getPosts() (*[]post_doc, error) {
	var result []post_doc
	cursor, err := db.Database(DATABASE).Collection("posts").Find(
		ctx, bson.D{{"author", string(u.Username)}},
	)
	err = cursor.All(ctx, &result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
		return nil, errors.New("User does not yet exist")

	}
	return &result, nil

}

func getUserStat[T any](u *user_doc, collection string) (*[]T, error) {
	var result []T
	cursor, err := db.Database(DATABASE).Collection(collection).Find(
		ctx, bson.D{{"author", string(u.Username)}},
	)
	err = cursor.All(ctx, &result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
		return nil, errors.New("User does not yet exist")

	}
	return &result, nil
}

func (u *user_doc) Delete(db *mongo.Client, ctx context.Context) error {
	_, err := db.Database(DATABASE).Collection("users").DeleteOne(ctx, bson.D{{"username", string(u.Username)}})
	return err
}
func (u *user_doc) newPost(title string, content string, db *mongo.Client, ctx context.Context) (post_doc, error) {
	collection := db.Database(DATABASE).Collection("posts")
	p := post_doc{
		Title:          title,
		Content:        content,
		Author:         u.Username,
		DateOfCreation: time.Now(),
		Likes:          make([]Username, 0),
		Comments:       make([]comment_doc, 0),
	}
	_, err := collection.InsertOne(context.TODO(), p)

	return p, err
}

func (u *user_doc) newComment(title string, content string, db *mongo.Client, ctx context.Context) (comment_doc, error) {
	collection := db.Database(DATABASE).Collection("comments")
	c := comment_doc{
		Content:        content,
		Author:         u.Username,
		DateOfCreation: time.Now(),
		Likes:          make([]Username, 0),
	}
	_, err := collection.InsertOne(context.TODO(), c)
	return c, err
}

func (u *user_doc) newNotification(nt NotificationType, content string, db *mongo.Client, ctx context.Context) (notification_doc, error) {
	collection := db.Database(DATABASE).Collection("notifications")

	n := notification_doc{
		Type:    nt,
		Content: content,
		Author:  u.Username,
	}
	_, err := collection.InsertOne(context.TODO(), n)

	return n, err
}
