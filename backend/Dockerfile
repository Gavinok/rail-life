#build stage
FROM golang:alpine
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
EXPOSE 8000
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

CMD ["go","run","db.go", "endpoints.go", "messaging.go", "main.go", "userapi.go"]
