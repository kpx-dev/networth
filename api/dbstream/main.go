package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/api/lib"
)

var (
	plaid  = nwlib.NewPlaidClient()
	db     = nwlib.NewDynamoDBClient()
	snsARN = nwlib.GetEnv("SNS_TOPIC_ARN")
)

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	var msg string
	for _, record := range e.Records {
		if record.Change.StreamViewType != "NEW_IMAGE" {
			msg = fmt.Sprintf("Received %s. Not a NEW_IMAGE stream view type, ignoring.", record.Change.StreamViewType)
			log.Println(msg)
			nwlib.PublishSNS(snsARN, msg)
			return
		}

		switch record.EventName {
		case "INSERT", "MODIFY":
			username, tokens := tokens(record)
			transactions(username, tokens)
			// msg = fmt.Sprintf("Insert / modify event %s, %v", username, tokens)
			break
		case "REMOVE":
			username, tokens := tokens(record)
			msg = fmt.Sprintf("Remove event %s, %v", username, tokens)
			break
		default:
			msg = fmt.Sprintf("Unknown event %s", record.EventName)
		}

		log.Println(msg)
		nwlib.PublishSNS(snsARN, msg)
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
