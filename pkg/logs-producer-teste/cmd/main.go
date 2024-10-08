package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"logs_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	logEntry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "auth-service",
		Level:     "INFO",
		Message:   "User login successful",
		Context: map[string]interface{}{
			"userID": "12345",
		},
	}

	body, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Failed to marshal log entry: %s", err)
	}

	err = ch.Publish(
		"",
		"logs_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Println("Log published to RabbitMQ")
}
