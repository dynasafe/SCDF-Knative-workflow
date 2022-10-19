package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"knative-example/pkg"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ENV_LISTEN_PORT  = "LISTEN_PORT"
	ENV_REDIS_MASTER = "REDIS_MASTER"
	ENV_REDIS_HOST   = "REDIS_HOST"
	ENV_REDIS_USER   = "REDIS_USER"
	ENV_REDIS_PASS   = "REDIS_PASS"
	ENV_REDIS_DB     = "REDIS_DB"
)

type CmdEvent struct {
	Command string `json:"command"`
}

func handleRequest(rdb *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allow", http.StatusMethodNotAllowed)
			return
		}

		body := r.Body
		defer body.Close()

		content, err := ioutil.ReadAll(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read message: %v", err), http.StatusInternalServerError)
			return
		}

		var event CmdEvent
		if err := json.Unmarshal(content, &event); err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshal event: %v", err), http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		taskId := fmt.Sprintf("%d", time.Now().Unix())
		fmt.Fprintf(w, "%s", taskId)

		key := fmt.Sprintf("SCDF-%s", taskId)
		req := pkg.SCDFRequest{
			Status:   false,
			TaskId:   taskId,
			ExitCode: 1,
			Command:  event.Command,
		}

		cmd := exec.Command("sh", "-c", event.Command)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Printf("Failed to execute this command: %v", err)
		} else {
			req.ExitCode = 0
			req.Status = true
		}
		fmt.Println(string(stdout))

		// send message
		msg, err := pkg.GenSCDFRequest(req)
		if err != nil {
			fmt.Printf("failed to generate request message: %v", err)
			return
		}
		rdsCmd := rdb.Set(ctx, key, msg, -1)
		if err := rdsCmd.Err(); err != nil {
			fmt.Printf("failed to set the value: %v", err)
			return
		}
	}
}

func main() {
	ops, err := pkg.GetRedisOptions(pkg.RedisConfig{
		Master:   os.Getenv(ENV_REDIS_MASTER),
		Host:     os.Getenv(ENV_REDIS_HOST),
		Username: os.Getenv(ENV_REDIS_USER),
		Password: os.Getenv(ENV_REDIS_PASS),
		Database: os.Getenv(ENV_REDIS_DB),
	})
	if err != nil {
		log.Fatalf("failed to get redis options: %v", err)
		return
	}

	rdb := redis.NewFailoverClient(ops)
	http.HandleFunc("/exec-job", handleRequest(rdb))
	port := "8084"
	if os.Getenv(ENV_LISTEN_PORT) != "" {
		port = os.Getenv(ENV_LISTEN_PORT)
	}
	panic(http.ListenAndServe(":"+port, nil))
}
