package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func decodeBody[T any](w http.ResponseWriter, r *http.Request) T {
	// Parse and decode the request body into a new `Credentials` instance
	var creds T
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error with request format")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	return creds

}
func Signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	// TODO ensure not signing in twice
	u, err := create_user(creds.Name,
		creds.Username,
		creds.Email,
		creds.DateOfBirth,
		db)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	log.Println("user added", u)
	json.NewEncoder(w).Encode(u)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	log.Println("user found", u)
	json.NewEncoder(w).Encode(u)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db)
	if err != nil {
		// Failed to read user
		log.Println("Error read user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	err = u.Delete(db)
	if err != nil {
		// Failed to delet user
		log.Println("Error delete user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	log.Println("user deleted", u)
	json.NewEncoder(w).Encode(u)
}

type new_post struct {
	U       user_doc `json:"u"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	log.Println("Creating new post")
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[new_post](w, r)
	u, err := read_user(creds.U.Username, db)
	if err != nil {
		log.Println(err)
	}
	post, err := u.newPost(creds.Title, creds.Content, db)
	if err != nil {
		log.Println(err)
	}
	log.Println("user post created", post)
	json.NewEncoder(w).Encode(post)
}

type new_friend struct {
	U          user_doc `json:"u"`
	FriendName string   `json:"friend_name"`
}

func NewFriend(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	log.Println("Adding a new friend")
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[new_friend](w, r)
	u, err := read_user(creds.U.Username, db)
	if err != nil {
		log.Println(err)
	}
	me, err := u.newFriend(creds.FriendName, db)
	if err != nil {
		log.Println("Error in NewFriend ", err)
	}
	log.Println("user as new friend", creds.FriendName)
	json.NewEncoder(w).Encode(me)
}

type new_comment struct {
	U            user_doc `json:"u"`
	ArticleTitle string   `json:"article_title"`
	Content      string   `json:"content"`
}

func NewComment(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	// TODO need a way to lookup a post to COMMENT on
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[new_comment](w, r)

	u, err := read_user(creds.U.Username, db)
	comment, err := u.newComment(creds.ArticleTitle, creds.Content, db)
	if err != nil {
		log.Println(err)
	}
	log.Println("user comment created", comment)
	json.NewEncoder(w).Encode(comment)
}

// Example of receiving realtime notifications
func ReadNotifications(w http.ResponseWriter, r *http.Request) {
	if mqconn == nil {
		mqconn, _ = rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
	}
	if (*r).Method == "OPTIONS" {
		return
	}
	// TODO need a way to lookup a post to COMMENT on
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db)
	if err != nil {
		log.Println("Error failed to find user in DB")
		w.WriteHeader(http.StatusBadRequest)
	}
	if u == nil {
		log.Println("Error user is nil")
		w.WriteHeader(http.StatusBadRequest)
	}
	notifiactions := u.getNotifications()
	json.NewEncoder(w).Encode(notifiactions)
}

// Example of sending realtime notifications
func SendNotifications(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	if mqconn == nil {
		mqconn, _ = rabbitMQDial("amqp://guest:guest@" + mqHOST + ":5672/")
	}
	// TODO need a way to lookup a post to COMMENT on
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db)
	if err != nil {
		log.Println("Error failed to find user in DB")
		w.WriteHeader(http.StatusBadRequest)
	}
	if u == nil {
		log.Println("Error user is nil")
		w.WriteHeader(http.StatusBadRequest)
	}
	ch := connectToQueue(mqconn)
	log.Println("notification queue connected")
	for i := 0; i < 11; i++ {
		ch.postNotification(u, "hello")
		time.Sleep(1 * time.Second)
	}
	log.Println("notifications sent")
	w.WriteHeader(http.StatusOK)

}
func AccountStats(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse and decode the request body into a new `Credentials` instance
	creds := &user_doc{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error with request format")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	u, err := read_user(creds.Username, db)

	if err != nil {
		log.Println(err)
	}
	stats, err := u.getStats()
	if err != nil {
		log.Println(err)
	}
	log.Println("stats found", stats)
	json.NewEncoder(w).Encode(stats)
}
