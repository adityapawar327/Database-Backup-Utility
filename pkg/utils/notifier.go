package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func SendSlackNotification(webhookURL string, message string) error {
	if webhookURL == "" {
		return nil // No webhook URL provided, skip notification
	}

	payload := SlackMessage{Text: message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send slack request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned status code: %d", resp.StatusCode)
	}

	return nil
}
