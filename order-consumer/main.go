package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	esTopic := os.Getenv("KAFKA_TOPIC")
	esURL := os.Getenv("ES_URL")
	groupID := os.Getenv("GROUP_ID")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   esTopic,
		GroupID: groupID,
	})
	defer reader.Close()

	fmt.Println("Consumer started, waiting for messages...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		var event map[string]interface{}
		json.Unmarshal(msg.Value, &event)

		payload, ok := event["payload"].(map[string]interface{})
		if !ok {
			continue
		}

		after, ok := payload["after"].(map[string]interface{})
		if !ok {
			continue
		}

		id := after["id"]
		docURL := fmt.Sprintf("%s/orders/_doc/%.0f", esURL, id)
		body, _ := json.Marshal(after)

		req, _ := http.NewRequest("PUT", docURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("ES error:", err)
			continue
		}
		resp.Body.Close()

		fmt.Printf("Indexed order %v to ES at %s\n", id, time.Now().Format("15:04:05"))
	}
}
