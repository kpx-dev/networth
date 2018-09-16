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
	plaidClient    = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	kms            = nwlib.NewKMSClient()
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
)

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	// TODO: https://github.com/aws/aws-lambda-go/issues/58

	for _, record := range e.Records {
		// nwlib.PublishSNS(snsARN, "record.Change.StreamViewType "+record.Change.StreamViewType)
		// nwlib.PublishSNS(snsARN, fmt.Sprintf("raw record %+v", record))
		fmt.Printf("raw record %+v", record)

		switch record.EventName {
		case "INSERT", "MODIFY":
			key := record.Change.Keys["key"].String()
			username := strings.Split(key, ":")[0]
			sort := record.Change.Keys["sort"].String()

			// each user have 2 sort keys for token: all, ins_XXX
			if strings.HasSuffix(key, ":token") && strings.HasPrefix(sort, "ins_") {
				tokens := record.Change.NewImage["tokens"].List()
				nwlib.PublishSNS(snsARN, fmt.Sprintf("len(newTokens): %d", len(tokens)))

				newToken := tokens[len(tokens)-1].Map()

				nwlib.PublishSNS(snsARN, fmt.Sprintf("len(record.Change.OldImage): %d", len(record.Change.OldImage)))
				if len(record.Change.OldImage) > 0 {
					go appendToken(username, newToken)
				}

				accessToken, err := kms.Decrypt(newToken["access_token"].String())

				if err != nil {
					return
				}

				go syncTransactions(username, accessToken)
				go syncAccounts(username, sort, accessToken)
			} else if strings.HasSuffix(key, ":account") {
				if sort == nwlib.DefaultSortValue {
					go syncNetworth(username)
				} else if strings.HasPrefix(sort, "ins_") && len(record.Change.OldImage) > 0 {
					// each user has 2 keys for account: all, ins_XXX
					nwlib.PublishSNS(snsARN, "about to sync appendAccount...")
					accounts := record.Change.NewImage["accounts"].List()
					go appendAccount(username, accounts)
				}
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
