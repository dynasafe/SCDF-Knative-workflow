package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var (
	ENV_KAFKA_HOST  = "KAFKA_HOST"
	ENV_KAFKA_TOPIC = "KAFKA_TOPIC"
	ENV_LISTEN_PORT = "LISTEN_PORT"
)

func getKafkaReader(kafkaURL, topic string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func readKafkaMessage(conn *websocket.Conn) error {
	kr := getKafkaReader(os.Getenv(ENV_KAFKA_HOST), os.Getenv(ENV_KAFKA_TOPIC))
	defer kr.Close()

	defer conn.Close()
	for {
		m, err := kr.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, m.Value); err != nil {
			log.Fatalf("failed to write message: %v", err)
		}
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go readKafkaMessage(conn)
}

func main() {
	log.Println("starting to create a new socket server")
	http.HandleFunc("/ws", handleRequest)
	port := "8081"
	if os.Getenv(ENV_LISTEN_PORT) != "" {
		port = os.Getenv(ENV_LISTEN_PORT)
	}
	panic(http.ListenAndServe(":"+port, nil))
}
