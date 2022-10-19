package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

/*
	REDIS_MASTER=mymaster \
	REDIS_HOST=10.250.75.117:26379 \
	REDIS_PASS=str0ng_passw0rd \
	go run write-redis/main.go
*/

func main() {
	redisOptions := &redis.FailoverOptions{
		MasterName:    os.Getenv("REDIS_MASTER"),
		SentinelAddrs: []string{os.Getenv("REDIS_HOST")},
	}

	// password
	if os.Getenv("REDIS_PASS") != "" {
		redisOptions.Password = os.Getenv("REDIS_PASS")
	}

	// database
	if os.Getenv("REDIS_DB") != "" {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			fmt.Printf("invalid database: %v", err)
			return
		}
		redisOptions.DB = db
	}

	// username
	if os.Getenv("REDIS_USER") != "" {
		redisOptions.Username = os.Getenv("REDIS_USER")
	}

	rdb := redis.NewFailoverClient(redisOptions)
	ctx := context.Background()
	cmd := rdb.Set(ctx, "test", "xxx", -1)
	if err := cmd.Err(); err != nil {
		fmt.Printf("failed to write this record: %v", err)
		return
	}
	fmt.Println("OK")
}
