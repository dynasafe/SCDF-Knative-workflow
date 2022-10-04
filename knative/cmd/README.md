## Kafka Client

### Overview
```
A simple tool to invoke SCDF's API

Usage:
  SCDF-Client [flags]

Flags:
  -h, --help                   help for SCDF-Client
      --kafka-host string      Please type Kafka endpoint
      --kafka-topic string     Please type Kafka topic
      --sasl-password string   Please type SASL password
      --sasl-username string   Please type SASL username
      --scdf-host string       Please type SCDF endpoint
      --skip-tls               insecure-skip-tls-verify
```

### Usage

type the command like the following:


kafka mode:

```
    client \
      --mode=kafka \
			--kafka-host=192.168.50.135:9092 \
			--kafka-topic=quickstart-events \
			--scdf-host="http://dataflow.prd.tanzu/tasks/executions?name=java-task03" \
			--sasl-username=client \
			--sasl-password=client-secret \
			--skip-tls=false
```

redis mode:

```
    client \
      --mode=redis \
      --scdf-host="http://dataflow.prd.tanzu/tasks/executions?name=java-task03" \
      --redis-host=10.250.75.117:26379 \
      --redis-db=0 \
      --redis-master-name=mymaster \
      --redis-pass=str0ng_passw0rd
      --redis-username=default
```