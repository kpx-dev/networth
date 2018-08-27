package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// func handler(ctx context.Context, snsEvent events.SNSEvent) {
// 	messages := ""
// 	for _, record := range snsEvent.Records {
// 		messages += record.SNS.Message + "\n"
// 	}

func handleAlert() {
	messages := "test"
	payload := []byte(`{"text":"` + messages + `"}`)
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")

	if _, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payload)); err != nil {
		fmt.Println("Problem with Slack", err)
	}
}
