package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SCDFRequest struct {
	Status   bool   `json:"status"`
	TaskId   string `json:"taskid"`
	ExitCode int    `json:"exit"`
}

func CallSCDFAPI(scdfHost string) (string, error) {
	fmt.Println("call SCDF API")
	resp, err := http.Post(scdfHost, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to call SCDF API: %v", err)
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return "", fmt.Errorf("failed to get the correct code")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get response body")
	}

	taskId := string(body)
	return taskId, nil
}

func ParseSCDFRequest(val []byte) (SCDFRequest, error) {
	var req SCDFRequest
	if err := json.Unmarshal(val, &req); err != nil {
		return SCDFRequest{}, fmt.Errorf("error request format: %s", string(val))
	}
	return req, nil
}
