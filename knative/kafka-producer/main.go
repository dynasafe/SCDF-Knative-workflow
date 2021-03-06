package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	ENV_KAFKA_HOST  = "KAFKA_HOST"
	ENV_KAFKA_TOPIC = "KAFKA_TOPIC"
	ENV_LISTEN_PORT = "LISTEN_PORT"
	ENV_USERNAME    = "SASL_USERNAME"
	ENV_PASSWORD    = "SASL_PASSWORD"
	ENV_SKIP_TLS    = "SKIP_TLS"
)

type SCDFEvent struct {
	Message string `json:"message"`
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv(ENV_USERNAME), os.Getenv(ENV_PASSWORD))
	if err != nil {
		panic(err)
	}

	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	if os.Getenv(ENV_SKIP_TLS) == "true" {
		sharedTransport.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &kafka.Writer{
		Addr:      kafka.TCP(kafkaURL),
		Topic:     topic,
		Transport: sharedTransport,
		Balancer:  &kafka.LeastBytes{},
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
