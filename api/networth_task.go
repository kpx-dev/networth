package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func networthTaskHandler(ctx context.Context, scheduleEvent events.CloudWatchEvent) {
	// messages := ""
	// for _, record := range snsEvent.Records {
	// 	messages += record.SNS.Message + "\n"
	// }

	// messages := "test"
	// payload := []byte(`{"text":"` + messages + `"}`)
	// slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")

	// if _, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payload)); err != nil {
	// 	fmt.Println("Problem with Slack", err)
	// }
}
