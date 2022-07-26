package main

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	ENV_KAFKA_HOST  = "KAFKA_HOST"
	ENV_KAFKA_TOPIC = "KAFKA_TOPIC"
	ENV_USERNAME    = "SAML_USERNAME"
	ENV_PASSWORD    = "SAML_PASSWORD"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv(ENV_USERNAME), os.Getenv(ENV_PASSWORD))
	if err != nil {
		panic(err)
	}

	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}

	return &kafka.Writer{
		Addr:      kafka.TCP(kafkaURL),
		Topic:     topic,
		Transport: sharedTransport,
		Balancer:  &kafka.LeastBytes{},
	}
}

func main() {
	kafkaURL := os.Getenv(ENV_KAFKA_HOST)
	topic := os.Getenv(ENV_KAFKA_TOPIC)
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	if err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("asd"),
	}); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
