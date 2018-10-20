package nwlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SlackBody struct
type SlackBody struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Channel  string `json:"channel"`
}

var cfg = LoadAWSConfig()
var snsClient = sns.New(cfg)

// PublishSNS publish message to SNS topic
func PublishSNS(arn string, message string) error {
	if message == "" {
		return errors.New("Cannot publish SNS, empty message")
	}

	request := snsClient.PublishRequest(&sns.PublishInput{
		Message:  &message,
		TopicArn: &arn,
	})

	_, err := request.Send()
	if err != nil {
		log.Printf("Problem publishing to SNS topic: %+v", err)
	}

	return err
}

// PublishSlack publish message to Slack
func PublishSlack(webhook string, message string, channel string) error {
	if message == "" {
		return errors.New("Cannot publish Slack, empty message")
	}

	slackBody := &SlackBody{
		Text: message,
	}

	if channel != "" {
		slackBody.Channel = channel
	}

	payload, err := json.Marshal(slackBody)

	if err != nil {
		log.Printf("Problem converting Slack body to json: %+v\n", err)
		return err
	}

	if res, err := http.Post(webhook, "application/json", bytes.NewBuffer(payload)); err != nil || res.StatusCode != http.StatusOK {
		log.Printf("Problem sending message using Slack: %+v", err)
		return err
	}

	return err
}
