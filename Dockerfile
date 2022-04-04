FROM golang:1.17.5-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o server ./cmd/billingapp
