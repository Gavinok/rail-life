#!/bin/sh

ENDPOINT="http://localhost:80"

# Create a user and a new post
curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff", "Password": "hello" }' "$ENDPOINT/signup"

curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/signin"

# Signup a user to use as a friend
curl -X POST -d '{ "Name": "chad", "Username": "chadman", "Email" : "chad@gmail.com",  "DateOfBirth": "1997ff" }' "$ENDPOINT/signup"
curl -X POST -d '{ "u": { "Username": "gavman" } , "friend_name": "chadman"}' "$ENDPOINT/newfriend"

# Create 10 posts to create 10 notifications for chad
bash ./create_ten.sh
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/stats"

curl -X POST -d '{ "u": { "Name": "chad", "Username": "chadman"} , "title": "I am chad", "content": "yo it is chad" }' "$ENDPOINT/post"
curl -X POST -d '{ "u": { "Username": "gavman" } , "article_title": "I am chad", "content": "sup chad" }' "$ENDPOINT/comment"

# get the stats for chad
curl -X POST -d '{ "Username": "chadman"}' "$ENDPOINT/stats"

echo "get the 10 most recent notifications"
curl -X POST -d '{ "Username": "chadman"}' "$ENDPOINT/notifications"

curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff" }' "$ENDPOINT/signup"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/signin"

# delete user
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/delete"
