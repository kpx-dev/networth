package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

// var (
// 	snsTopicARN     = GetEnv("SNS_TOPIC_ARN")
// 	slackWebhookURL = GetEnv("SLACK_WEBHOOK_URL")
// 	slackChannel    = GetEnv("SLACK_CHANNEL")
// )

func TestAccounts(t *testing.T) {
	username := "test@networth.app"
	token := "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"

	if err := syncAccounts(username, token); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}
