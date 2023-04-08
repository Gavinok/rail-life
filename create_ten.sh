#!/bin/sh

ENDPOINT="http://localhost:80"

# Create 10 posts
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello1  world", "content": "what is up people1 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello2  world", "content": "what is up people2 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello3  world", "content": "what is up people3 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello4  world", "content": "what is up people4 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello5  world", "content": "what is up people5 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello6  world", "content": "what is up people6 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello7  world", "content": "what is up people7 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello8  world", "content": "what is up people8 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello9  world", "content": "what is up people9 " }' "$ENDPOINT/post"
curl -X POST -d '{ "u": {"Username": "gavman" } , "title": "hello10  world", "content": "what is up people10 " }' "$ENDPOINT/post"
curl -X POST -d '{ "Username": "gavman"}' "$ENDPOINT/stats"
