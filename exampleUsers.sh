#!/bin/sh

ENDPOINT="http://localhost:8000"

# Create a user and a new post
echo "should have cache miss"
curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff", "Password": "hello" }' "$ENDPOINT/signup"

echo "should have cache hit"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/signin"

echo "should have cache hit"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/signin"

echo "crate a new user and add them as a friend"
curl -X POST -d '{ "Name": "chad", "Username": "chadman", "Email" : "chad@gmail.com",  "DateOfBirth": "1997ff" }' "$ENDPOINT/signup"
curl -X POST -d '{ "u": { "Username": "gavman" } , "friend_name": "chadman"}' "$ENDPOINT/newfriend"

echo "should have cache hit"
curl -X POST -d '{ "u": { "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff" } , "title": "hello world", "content": "what is up people" }' "$ENDPOINT/post"
# curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/comment"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/stats"

curl -X POST -d '{ "u": { "Name": "chad", "Username": "chadman"} , "title": "I am chad", "content": "yo it is chad" }' "$ENDPOINT/post"
curl -X POST -d '{ "u": { "Username": "gavman" } , "article_title": "I am chad", "content": "sup chad" }' "$ENDPOINT/comment"

# echo "Force notifications for this user"
# curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff" }' "$ENDPOINT/forceNotifications" &

echo "get the 10 most recent notifications"
curl -X POST -d '{ "Username": "chadman"}' "$ENDPOINT/notifications"

curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff" }' "$ENDPOINT/signup"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/signin"

# delete user
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/delete"
