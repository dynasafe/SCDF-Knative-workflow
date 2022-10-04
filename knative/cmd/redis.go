package main

import (
	"context"
	"fmt"
	"knative-example/pkg"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	REDIS_CONNECTION_NAME = "redis"
)

func newRedisConn() ConnectionType {
	return &RedisConn{}
}

type RedisConn struct{}

func (RedisConn) Name() string {
	return REDIS_CONNECTION_NAME
}

type SCDFSignal struct {
	Request pkg.SCDFRequest
	Error   error
}

/*
	go run ./cmd/...
		--mode=redis \
		--scdf-host="http://dataflow.prd.tanzu/tasks/executions?name=java-task03" \
		--redis-host=10.250.75.117:26379 \
		--redis-db=0 \
		--redis-master-name=mymaster \
		--redis-pass=str0ng_passw0rd
		--redis-username=default
*/
func (RedisConn) Connect() error {
	redisOptions := &redis.FailoverOptions{
		MasterName:    redisMaster,
		SentinelAddrs: []string{redisHost},
		Password:      redisPass,
	}

	// database
	if redisDB != "" {
		db, err := strconv.Atoi(redisDB)
		if err != nil {
			return fmt.Errorf("invalid database: %v", err)
		}
		redisOptions.DB = db
	}

	// username
	if redisUsername != "" {
		redisOptions.Username = redisUsername
	}

	rdb := redis.NewFailoverClient(redisOptions)
	ctx := context.Background()
	taskId, err := pkg.CallSCDFAPI(scdfHost)
	if err != nil {
		return fmt.Errorf("failed to invoke SCDF API: %v", err)
	}

	ch := make(chan SCDFSignal, 1)
	defer close(ch)
	go func(red *redis.Client, ch chan SCDFSignal, taskId string) {
		var (
			getKeyErr error
			request   pkg.SCDFRequest
			count     int
		)
		keyName := fmt.Sprintf("SCDF-%s", taskId)
		fmt.Printf("\n[%s] start reading... \n\n", keyName)
		for {
			if count != 0 {
				time.Sleep(10 * time.Second)
			}
			val, err := rdb.Get(ctx, keyName).Result()
			if err == redis.Nil {
				count++
				fmt.Printf("[%d] key does not exist \n", count)
				continue // if value doesn't exist, keep searching
			} else if err != nil {
				getKeyErr = fmt.Errorf("failed to get this key: %v", err)
				break
			}

			req, err := pkg.ParseSCDFRequest([]byte(val))
			if err != nil {
				getKeyErr = fmt.Errorf("failed to parse SCDF request: %v", err)
				break
			}
			request = req
			break
		}

		if getKeyErr != nil {
			ch <- SCDFSignal{Error: getKeyErr}
			return
		}
		ch <- SCDFSignal{Request: request}
	}(rdb, ch, taskId)

	signal := <-ch
	if signal.Error != nil {
		return fmt.Errorf("failed to get the value: %v", signal.Error)
	}

	exitCode := signal.Request.ExitCode
	if exitCode != 0 {
		return fmt.Errorf("failed to execute the job: exit(%d)", exitCode)
	}

	fmt.Printf("batch job has done, status: %t\n", signal.Request.Status)
	return nil
}
