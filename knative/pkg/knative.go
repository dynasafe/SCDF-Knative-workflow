package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type KnativeRequest struct {
	Endpoint string
	Host     string
	Command  string
}

func CallKnativeAPI(kreq KnativeRequest) (string, error) {
	jsonBody, err := json.Marshal(struct {
		Command string `json:"command"`
	}{
		Command: kreq.Command,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, kreq.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to new a request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Host = kreq.Host

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call API: %v", err)
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return "", fmt.Errorf("failed to get the correct code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get response body")
	}

	taskId := string(body)
	return taskId, nil
}
