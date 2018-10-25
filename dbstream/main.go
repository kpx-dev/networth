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

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		partitionKey := record.Change.Keys["id"].String()

		switch record.EventName {
		case "INSERT", "MODIFY":
			if strings.HasSuffix(partitionKey, ":webhook") {
				if err := handleInsertModifyWebhook(record); err != nil {
					log.Printf("Problem insert / modify webhook: %+v\n", err)
				}
			} else if strings.HasSuffix(partitionKey, ":token") {
				if err := handleInsertModifyToken(record); err != nil {
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
