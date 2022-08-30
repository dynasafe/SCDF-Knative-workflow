package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"knative-example/pkg"

	"github.com/spf13/cobra"
)

var (
	kafkaHost    string
	kafkaTopic   string
	scdfHost     string
	saslUsername string
	saslPassword string
	skipTLS      bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&kafkaHost, "kafka-host", "", "Please type Kafka endpoint")
	rootCmd.PersistentFlags().StringVar(&kafkaTopic, "kafka-topic", "", "Please type Kafka topic")
	rootCmd.PersistentFlags().StringVar(&scdfHost, "scdf-host", "", "Please type SCDF endpoint")
	rootCmd.PersistentFlags().StringVar(&saslUsername, "sasl-username", "", "Please type SASL username")
	rootCmd.PersistentFlags().StringVar(&saslPassword, "sasl-password", "", "Please type SASL password")
	rootCmd.PersistentFlags().BoolVar(&skipTLS, "skip-tls", false, "insecure-skip-tls-verify")
}

type SCDFRequest struct {
	Status bool   `json:"status"`
	TaskId string `json:"taskid"`
}

// usage:
// 		go run cmd/main.go \
//			--kafka-host=192.168.50.135:9092 \
//			--kafka-topic=quickstart-events \
//			--scdf-host="http://dataflow.prd.tanzu/tasks/executions?name=java-task03" \
//			--sasl-username=client \
//			--sasl-password=client-secret \
//			--skip-tls=false
var rootCmd = &cobra.Command{
	Use:   "SCDF-Client",
	Short: "A simple tool to invoke SCDF's API",
	Long:  "A simple tool to invoke SCDF's API",
	RunE: func(cmd *cobra.Command, args []string) error {
		kr := pkg.GetKafkaReader(pkg.KafkaRequest{
			KafkaURL: kafkaHost,
			Topic:    kafkaTopic,
			Username: saslUsername,
			Password: saslPassword,
			SkipTLS:  skipTLS,
		})
		defer kr.Close()

		fmt.Println("call SCDF API")
		resp, err := http.Post(scdfHost, "", nil)
		if err != nil {
			return fmt.Errorf("failed to call SCDF API: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to get request body")
		}

		taskId := string(body)
		fmt.Printf("[%s] start reading Kafka request", taskId)
		for {
			m, err := kr.ReadMessage(context.Background())
			if err != nil {
				log.Fatalf("failed to read message: %v", err)
				break
			}
			var req SCDFRequest
			if err := json.Unmarshal(m.Value, &req); err != nil {
				log.Fatalf("error request format: %v", err)
				log.Println(string(m.Value))
			}
			if req.TaskId == taskId {
				fmt.Printf("batch job has done, status: %t", req.Status)
				break
			}
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
