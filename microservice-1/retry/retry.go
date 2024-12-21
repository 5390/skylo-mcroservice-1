package retry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"microservice-1/config"
	"net/http"
	"time"
)

type RetryHandler struct {
	url        string
	retryDelay time.Duration
}

func NewRetryHandler(config config.RetryConfig) *RetryHandler {
	return &RetryHandler{
		url:        config.TargetURL,
		retryDelay: config.RetryDelay,
	}
}

func (r *RetryHandler) ProcessMessage(message string) error {
	for {
		err := r.sendToMicroservice2(message)
		if err != nil {
			log.Printf("Retrying in %v seconds. Error: %v\n", r.retryDelay.Seconds(), err)
			time.Sleep(r.retryDelay)
			continue
		}
		return nil
	}
}

func (r *RetryHandler) sendToMicroservice2(message string) error {
	// Create a map to hold the JSON structure
	payload := map[string]string{
		"data": message,
	}

	// Marshal the map into a JSON byte slice
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create a new POST request with the marshaled JSON data
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		return errors.New("non-200 response from Microservice-2")
	}

	return nil
}
