package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	ENV_KAFKA_HOST  = "KAFKA_HOST"
	ENV_KAFKA_TOPIC = "KAFKA_TOPIC"
	ENV_USERNAME    = "SAML_USERNAME"
	ENV_PASSWORD    = "SAML_PASSWORD"
)

func newKafkaReader(kafkaURL, topic string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv(ENV_USERNAME), os.Getenv(ENV_PASSWORD))
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		Dialer:  dialer,
	})
}

func main() {
	kafkaURL := os.Getenv(ENV_KAFKA_HOST)
	topic := os.Getenv(ENV_KAFKA_TOPIC)
	kr := newKafkaReader(kafkaURL, topic)
	defer kr.Close()

	for {
		m, err := kr.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
			break
		}
		fmt.Println(m.Value)
	}
	os.Exit(0)
}
