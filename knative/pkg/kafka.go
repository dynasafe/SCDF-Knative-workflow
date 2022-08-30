package pkg

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaRequest struct {
	KafkaURL string
	Topic    string
	Username string
	Password string
	SkipTLS  bool
}

func GetKafkaReader(req KafkaRequest) *kafka.Reader {
	brokers := strings.Split(req.KafkaURL, ",")
	mechanism, err := scram.Mechanism(scram.SHA256, req.Username, req.Password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	if req.SkipTLS {
		dialer.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   req.Topic,
		Dialer:  dialer,
	})
}
