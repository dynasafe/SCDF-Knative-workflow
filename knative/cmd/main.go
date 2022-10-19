package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type ConnectionType interface {
	Name() string
	Connect() error
}

var (
	redisHost       string
	redisDB         string
	redisMaster     string
	redisUsername   string
	redisPass       string
	connMode        string
	kafkaHost       string
	kafkaTopic      string
	scdfHost        string
	saslUsername    string
	saslPassword    string
	skipTLS         bool
	connectionTypes []ConnectionType
	method          string
	knHost          string
	knEndpoint      string
	knCommand       string
)

func init() {
	registerTypes(
		newKafkaConn(),
		newRedisConn(),
	)
	rootCmd.PersistentFlags().StringVar(&connMode, "mode", "", fmt.Sprintf("Please type Connection mode: (ie: %s)", getModeNames()))
	rootCmd.PersistentFlags().StringVar(&method, "method", "scdf", "Please type executing method: (ie: scdf or knative)")

	// redis
	rootCmd.PersistentFlags().StringVar(&knHost, "kn-host", "", "Please type Knative Host")
	rootCmd.PersistentFlags().StringVar(&knEndpoint, "kn-endpoint", "", "Please type Knative endpoint")
	rootCmd.PersistentFlags().StringVar(&knCommand, "kn-cmd", "", "Please input the command")
	rootCmd.PersistentFlags().StringVar(&redisHost, "redis-host", "", "Please type Redis endpoint")
	rootCmd.PersistentFlags().StringVar(&redisDB, "redis-db", "", "[option] Please type Redis DB")
	rootCmd.PersistentFlags().StringVar(&redisMaster, "redis-master-name", "", "Please type Redis Sentinel master name")
	rootCmd.PersistentFlags().StringVar(&redisUsername, "redis-username", "", "[option] Please type Redis username")
	rootCmd.PersistentFlags().StringVar(&redisPass, "redis-pass", "", "Please type Redis password")

	// kafka
	rootCmd.PersistentFlags().StringVar(&kafkaHost, "kafka-host", "", "Please type Kafka endpoint")
	rootCmd.PersistentFlags().StringVar(&kafkaTopic, "kafka-topic", "", "Please type Kafka topic")
	rootCmd.PersistentFlags().StringVar(&scdfHost, "scdf-host", "", "Please type SCDF endpoint")
	rootCmd.PersistentFlags().StringVar(&saslUsername, "sasl-username", "", "Please type SASL username")
	rootCmd.PersistentFlags().StringVar(&saslPassword, "sasl-password", "", "Please type SASL password")
	rootCmd.PersistentFlags().BoolVar(&skipTLS, "skip-tls", false, "insecure-skip-tls-verify")
}

func registerTypes(types ...ConnectionType) {
	connectionTypes = append(connectionTypes, types...)
}

var rootCmd = &cobra.Command{
	Use:   "SCDF-Client",
	Short: "A simple tool to invoke SCDF's API",
	Long:  "A simple tool to invoke SCDF's API",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, conn := range connectionTypes {
			if connMode == conn.Name() {
				return conn.Connect()
			}
		}
		return fmt.Errorf("invalid connection type: %s", connMode)
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

func getModeNames() string {
	names := make([]string, 0, len(connectionTypes))
	for _, conn := range connectionTypes {
		names = append(names, conn.Name())
	}
	return strings.Join(names[:], ",")
}
