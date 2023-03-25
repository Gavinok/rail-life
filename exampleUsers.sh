#!/bin/sh

# Create a user and a new post

curl -X POST -d '{ "Name": "gavin", "Username": "gavman", "Email" : "gavinfreeborn@gmail.com",  "DateOfBirth": "1997ff" }' "http://localhost:8000/signup"
curl -X POST -d '{ "Username": "gavman"}' "http://localhost:8000/signin"
curl -X POST -d '{ "Username": "gavman"}' "http://localhost:8000/delete"
curl -X POST -d '{ "Username": "gavman"}' "http://localhost:8000/post"
curl -X POST -d '{ "Username": "gavman"}' "http://localhost:8000/comment"
curl -X POST -d '{ "Username": "gavman"}' "http://localhost:8000/stats"

