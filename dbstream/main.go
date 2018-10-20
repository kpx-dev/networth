package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/lib"
)

var (
	plaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv       = nwlib.GetEnv("PLAID_ENV")
	plaidClient    = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	kms            = nwlib.NewKMSClient()
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
)

func extractCompositeKeys(record events.DynamoDBEventRecord) (string, string, string) {
	partitionKey := record.Change.Keys["id"].String()
	sortKey := record.Change.Keys["sort"].String()
	username := strings.Split(partitionKey, ":")[0]

	return partitionKey, sortKey, username
}

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		key, sort, username := extractCompositeKeys(record)

		switch record.EventName {
		case "INSERT", "MODIFY":
			if key == "webhook" {
				if err := handleInsertModifyWebhook(record); err != nil {
					log.Printf("Problem insert / modify webhook: %+v\n", err)
				}
			} else if strings.HasSuffix(key, ":token") {
				if err := handleInsertModifyToken(username, sort, record); err != nil {
					log.Printf("Problem insert / modify token: %+v\n", err)
				}
			}
			break
		default:
			log.Printf("DynamoDB stream unknown event: %s, records: %+v\n", record.EventName, record)
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
