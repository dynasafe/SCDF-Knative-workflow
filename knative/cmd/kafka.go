package main

import (
	"context"
	"fmt"
	"knative-example/pkg"
)

/* usage:
 		go run cmd/main.go \
			--mode=kafka \
			--kafka-host=192.168.50.135:9092 \
			--kafka-topic=quickstart-events \
			--scdf-host="http://dataflow.prd.tanzu/tasks/executions?name=java-task03" \
			--sasl-username=client \
			--sasl-password=client-secret \
			--skip-tls=false
*/
const (
	KAFKA_CONNECTION_NAME = "kafka"
)

type KafkaConn struct{}

func newKafkaConn() ConnectionType {
	return &KafkaConn{}
}

func (KafkaConn) Name() string {
	return KAFKA_CONNECTION_NAME
}

func (KafkaConn) Connect() error {
	kr := pkg.GetKafkaReader(pkg.KafkaRequest{
		KafkaURL: kafkaHost,
		Topic:    kafkaTopic,
		Username: saslUsername,
		Password: saslPassword,
		SkipTLS:  skipTLS,
	})
	defer kr.Close()
	taskId, err := pkg.CallSCDFAPI(scdfHost)
	if err != nil {
		return fmt.Errorf("failed to invoke SCDF API: %v", err)
	}
	fmt.Printf("[%s] start reading Kafka request", taskId)
	var exitCode int
	for {
		m, err := kr.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to read message: %v", err)
		}
		req, err := pkg.ParseSCDFRequest(m.Value)
		if err != nil {
			return fmt.Errorf("failed to parse the request")
		}
		if req.TaskId != taskId {
			continue
		}
		fmt.Printf("batch job has done, status: %t", req.Status)
		exitCode = req.ExitCode
		break
	}
	if exitCode != 0 {
		return fmt.Errorf("failed to execute the job: exit(%d)", exitCode)
	}
	return nil
}
