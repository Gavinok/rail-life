package main

import (
	"encoding/json"
	"log"
	"net/http"
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
		db,
		ctx)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	log.Println("user added", u)
}
func SignIn(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db, ctx)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	log.Println("user found", u)
}
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db, ctx)
	if err != nil {
		// Failed to read user
		log.Println("Error read user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	err = u.Delete(db, ctx)
	if err != nil {
		// Failed to delet user
		log.Println("Error delete user")
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	log.Println("user deleted", u)
}
func NewPost(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db, ctx)
	post, err := u.newPost("hello world", "Hey there", db, ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println("user post created", post)
}

func NewComment(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}
	// TODO need a way to lookup a post to COMMENT on
	// Parse and decode the request body into a new `Credentials` instance
	creds := decodeBody[user_doc](w, r)
	u, err := read_user(creds.Username, db, ctx)
	comment, err := u.newComment("hello world", "Hey there", db, ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println("user comment created", comment)
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
	u, err := read_user(creds.Username, db, ctx)
	posts, err := u.getPosts()
	comments, err := getUserStat[comment_doc](u, "comments")
	// TODO liked comments and posts
	if err != nil {
		log.Println(err)
	}
	log.Println("posts found", posts)
	log.Println("comments found", comments)
}
