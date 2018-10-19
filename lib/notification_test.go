package nwlib

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

var (
	snsTopicARN     = GetEnv("SNS_TOPIC_ARN")
	slackWebhookURL = GetEnv("SLACK_WEBHOOK_URL")
	slackChannel    = GetEnv("SLACK_CHANNEL")
)

func TestPublishSlack(t *testing.T) {
	// send to default channel
	if err := PublishSlack(slackWebhookURL, "test_publish_slack", slackChannel); err != nil {
		t.Errorf("Cannot publish Slack %v", err)
	}

	// send to custom channel
	if err := PublishSlack(slackWebhookURL, "test_publish_slack", "test"); err != nil {
		t.Errorf("Cannot publish Slack %v", err)
	}
}

func TestPublishSNS(t *testing.T) {
	if err := PublishSNS(snsTopicARN, "test_publish_sns"); err != nil {
		t.Errorf("Cannot publish SNS %v", err)
	}
}
