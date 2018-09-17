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
	wg             sync.WaitGroup
)

func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	// TODO: https://github.com/aws/aws-lambda-go/issues/58

	for _, record := range e.Records {
		switch record.EventName {
		case "INSERT", "MODIFY":
			key := record.Change.Keys["key"].String()
			username := strings.Split(key, ":")[0]
			sort := record.Change.Keys["sort"].String()

			// each user have 2 sort keys for token: all, ins_XXX
			if strings.HasSuffix(key, ":token") && strings.HasPrefix(sort, "ins_") {
				tokens := record.Change.NewImage["tokens"].List()
				newToken := tokens[len(tokens)-1].Map()

				if len(record.Change.OldImage) > 0 {
					appendToken(username, newToken)
				}

				accessToken, err := kms.Decrypt(newToken["access_token"].String())

				if err != nil {
					return
				}

				// TODO: make these into gorutines / wait group workers:
				// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
				// syncTransactions(username, accessToken)
				syncAccounts(username, sort, accessToken)
			} else if strings.HasSuffix(key, ":account") {
				if sort == nwlib.DefaultSortValue {
					syncNetworth(username)
				} else if strings.HasPrefix(sort, "ins_") && len(record.Change.OldImage) > 0 {
					// each user has 2 keys for account: all, ins_XXX
					// TODO: [bug] somehow have duplicate accounts in all key
					accounts := record.Change.NewImage["accounts"].List()
					appendAccount(username, accounts)
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
