package main

import (
	"context"
	"errors"
	"log"

	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDB() *dataStore {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://admin:admin@" + dbHOST + ":27017")
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		fmt.Println("Error connecting to DB")
		panic(err)
	}
	return &dataStore{
		client, ctx,
	}
}

type Username string
type user_doc struct {
	Name        string
	Username    Username
	Email       string
	DateOfBirth string
	Friends     []Username
	Password    string
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
type NotificationType int8

const (
	COMMENT_NOTIFICATION NotificationType = iota
	LIKE_NOTIFICATION
	POST_NOTIFICATION
)

type notification_doc struct {
	Type    NotificationType
	Content string
	Author  Username
}

const DATABASE = "social_media_app"

func read_user(user Username, db *dataStore) (*user_doc, error) {
	if db == nil {
		db = connectToDB()
	}

	// Check the cache for efficiency reasons
	usr, err := CheckCache(user_doc{Username: user}, mycache, cacheCTX)
	if err == nil {
		log.Println("Cache hit")
		return &usr, nil
	} else {
		log.Println("error with cache ", err)
	}

	var result user_doc
	err = db.db.Database(DATABASE).Collection("users").FindOne(
		db.ctx, bson.D{{"username", string(user)}},
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

func create_user(user string, username Username, email string, dateofbirth string, db *dataStore) (*user_doc, error) {
	u := user_doc{
		Name:        user,
		Username:    username,
		Email:       email,
		DateOfBirth: dateofbirth,
		Friends:     make([]Username, 0),
	}
	_, err := db.db.Database(DATABASE).Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		fmt.Println("Error adding user to db: ", err)
		panic(err)
	}

	err = Cache(u, mycache, cacheCTX)
	CheckCache(u, mycache, cacheCTX)
	if err != nil {
		fmt.Println("Error adding user to db: ", err)
		panic(err)
	}

	return &u, nil
}

func (u *user_doc) addFriend(friend Username, db *dataStore) (*user_doc, error) {
	upsert := true
	after := options.After
	var user_friend user_doc
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	ret := db.db.Database(DATABASE).Collection("users").FindOneAndUpdate(context.TODO(),
		bson.M{"username": string(u.Username)},
		bson.M{
			"$addToSet": bson.M{"friends": string(friend)},
		},
		&opt)
	if ret.Err() != nil {
		log.Println("error adding friend", ret.Err())
		return nil, ret.Err()
	}
	ret.Decode(&user_friend)
	return &user_friend, Cache(user_friend, mycache, cacheCTX)

}

type dataStore struct {
	db  *mongo.Client
	ctx context.Context
}

func (db *dataStore) Disconnect() {
	db.db.Disconnect(db.ctx)
}

func (db *dataStore) lookupPostsByAuthor(u *user_doc) (*[]post_doc, error) {
	var result []post_doc
	cursor, err := db.db.Database(DATABASE).Collection("posts").Find(
		db.ctx, bson.D{{"author", string(u.Username)}},
	)
	err = cursor.All(db.ctx, &result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
		return nil, errors.New("User does not yet exist")

	}
	return &result, err
}
func getUserStat[T any](u *user_doc, collection string) (*[]T, error) {
	var result []T
	cursor, err := db.db.Database(DATABASE).Collection(collection).Find(
		db.ctx, bson.D{{"author", string(u.Username)}},
	)
	err = cursor.All(db.ctx, &result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
		return nil, errors.New("User does not yet exist")

	}
	return &result, nil
}

func (db *dataStore) deleteUser(u *user_doc) error {
	_, err := db.db.Database(DATABASE).Collection("users").DeleteOne(db.ctx, bson.D{{"username", string(u.Username)}})
	return err
}
func (db *dataStore) storePost(p post_doc) error {
	collection := db.db.Database(DATABASE).Collection("posts")
	_, err := collection.InsertOne(context.TODO(), p)
	return err
}
func (db *dataStore) storeComment(post_title string, comment comment_doc) (*post_doc, error) {
	var result post_doc
	// 7) Create an instance of an options and set the desired options
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	ret := db.db.Database(DATABASE).Collection("posts").FindOneAndUpdate(context.TODO(),
		bson.M{"title": post_title},
		bson.M{
			"$addToSet": bson.M{"comments": comment},
		},
		&opt)
	if ret.Err() != nil {
		log.Println("error for comment insertion ", ret.Err())
		return nil, ret.Err()
	}
	ret.Decode(&result)
	return &result, nil
}

func (db *dataStore) storeNotification(u *user_doc, n notification_doc) (notification_doc, error) {
	collection := db.db.Database(DATABASE).Collection("notifications")
	_, err := collection.InsertOne(context.TODO(), n)

	return n, err
}
