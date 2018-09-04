package nwlib

import (
	"bytes"
	"log"
	"net/http"
)

var slackWebhookURL = GetEnv("SLACK_WEBHOOK_URL")

// Alert message to Slack. TODO: using sns...
func Alert(message string) {
	log.Println("Alerting to Slack: ", message)
	payload := []byte(`{"text":"` + message + `"}`)

	if _, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payload)); err != nil {
		log.Println("Cannot alert, problem with Slack", err)
	}
}
