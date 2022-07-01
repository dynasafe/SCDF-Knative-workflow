package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/segmentio/kafka-go"
)

var (
	ENV_KAFKA_HOST  = "KAFKA_HOST"
	ENV_KAFKA_TOPIC = "KAFKA_TOPIC"
	ENV_LISTEN_PORT = "LISTEN_PORT"
)

type SCDFEvent struct {
	Message string `json:"message"`
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func handleRequest(kw *kafka.Writer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allow", http.StatusMethodNotAllowed)
			return
		}

		body := r.Body
		defer body.Close()

		content, err := ioutil.ReadAll(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read message: %v", err), http.StatusInternalServerError)
			return
		}

		var event SCDFEvent
		if err := json.Unmarshal(content, &event); err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshal event: %v", err), http.StatusInternalServerError)
			return
		}

		if err := kw.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(event.Message),
		}); err != nil {
			http.Error(w, fmt.Sprintf("failed to send message: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {

	kafkaURL := os.Getenv(ENV_KAFKA_HOST)
	topic := os.Getenv(ENV_KAFKA_TOPIC)
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	http.HandleFunc("/send-msg", handleRequest(writer))
	port := "8084"
	if os.Getenv(ENV_LISTEN_PORT) != "" {
		port = os.Getenv(ENV_LISTEN_PORT)
	}
	panic(http.ListenAndServe(":"+port, nil))
}
