package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/api/lib"
)

var (
	plaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv       = nwlib.GetEnv("PLAID_ENV", "sandbox")
	kms            = nwlib.NewKMSClient()
	plaid          = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
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
			msg = fmt.Sprintf("DynamoDB stream: %s", record.EventName)
			key := record.Change.Keys["key"].String()
			username := strings.Split(key, ":")[0]
			sort := record.Change.Keys["sort"].String()

			if strings.HasSuffix(key, ":token") && strings.HasPrefix(sort, "ins_") {
				// TODO: https://github.com/aws/aws-lambda-go/issues/58
				tokens := record.Change.NewImage["tokens"].List()
				newToken := tokens[len(tokens)-1].Map()
				go appendToken(username, newToken)

				accessToken, err := kms.Decrypt(newToken["access_token"].String())

				if err != nil {
					return
				}

				go syncTransactions(username, accessToken)
				go syncAccounts(username, accessToken)
			}
			break
		default:
			msg = fmt.Sprintf("DynamoDB stream unknown event %s %+v", record.EventName, record)
		}

		log.Println(msg)
		nwlib.PublishSNS(snsARN, msg)
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
