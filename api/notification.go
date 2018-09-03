package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/networth-app/networth/api/lib"
)

var snsClient = sns.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
var snsTopicARN = nwlib.GetEnv("SNS_TOPIC_ARN", "")

// PublishSNS publish message to SNS topic
func PublishSNS(message string) {
	fmt.Println("Publishing message:", message)
	input := &sns.PublishInput{
		Message:  &message,
		TopicArn: &snsTopicARN,
	}
	_, err := snsClient.Publish(input)
	if err != nil {
		fmt.Println("Problem publishing to SNS topic", err)
	}
}
