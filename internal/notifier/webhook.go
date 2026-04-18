package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type webhookPayload struct {
	Text string `json:"text"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func sendWebhook(url, message string) error {
	if url == "" {
		return fmt.Errorf("notifier: webhook target URL is empty")
	}
	payload := webhookPayload{Text: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notifier: marshal payload: %w", err)
	}
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
