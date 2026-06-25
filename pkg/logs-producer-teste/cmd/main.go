package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	url := getenv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	queue := getenv("LOGS_QUEUE", "logs_queue")
	interval, _ := strconv.Atoi(getenv("PRODUCER_INTERVAL_SECONDS", "5"))
	if interval <= 0 {
		interval = 5
	}

	var conn *amqp.Connection
	var err error
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ indisponível (%v), tentando de novo...", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	levels := []string{"INFO", "WARNING", "ERROR"}
	i := 0
	for {
		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Service:   "demo-service",
			Level:     levels[i%len(levels)],
			Message:   "Evento sintético de demonstração",
			Context:   map[string]interface{}{"seq": i},
		}
		body, _ := json.Marshal(entry)
		if err := ch.Publish("", queue, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
			log.Printf("Failed to publish: %s", err)
		} else {
			log.Printf("Log publicado (seq=%d, level=%s)", i, entry.Level)
		}
		i++
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
