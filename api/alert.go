package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

func alert(message string) {
	log.Println("About to alert: ", message)
	payload := []byte(`{"text":"` + message + `"}`)
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")

	if _, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payload)); err != nil {
		log.Println("Cannot alert, problem with Slack", err)
	}
}
