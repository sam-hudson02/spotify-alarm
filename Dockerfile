FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY src src

RUN go build ./src/server/main.go

EXPOSE 8080
CMD ["./main"]
