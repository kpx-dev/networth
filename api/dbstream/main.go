package main

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/api/lib"
)

var (
	plaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv       = nwlib.GetEnv("PLAID_ENV", "sandbox")
	plaidClient    = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	kms            = nwlib.NewKMSClient()
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
	slackURL       = nwlib.GetEnv("SLACK_WEBHOOK_URL")
	wg             sync.WaitGroup
)

func extractCompositeKeys(record events.DynamoDBEventRecord) (string, string, string) {
	key := record.Change.Keys["id"].String()
	username := strings.Split(key, ":")[0]
	sort := record.Change.Keys["sort"].String()

	return key, username, sort
}

// TODO: https://github.com/aws/aws-lambda-go/issues/58
func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		key, username, sort := extractCompositeKeys(record)

		switch record.EventName {
		case "INSERT":
			if strings.HasSuffix(key, ":token") {
				handleInsertModifyToken(username, sort, record)
			} else if strings.HasSuffix(key, ":account") {
				handleInsertAccount(username, sort, record)
			}
			break

		case "MODIFY":
			if strings.HasSuffix(key, ":token") {
				handleInsertModifyToken(username, sort, record)
			} else if strings.HasSuffix(key, ":account") && sort == nwlib.DefaultSortValue {
				syncNetworth(username)
			}
			break

		default:
			log.Printf("DynamoDB stream unknown event %s %+v", record.EventName, record)
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
